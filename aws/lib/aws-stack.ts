import * as cdk from "aws-cdk-lib";
import { Construct } from "constructs";
import * as ecr from "aws-cdk-lib/aws-ecr";
import * as ecs from "aws-cdk-lib/aws-ecs";
import * as ec2 from "aws-cdk-lib/aws-ec2";
import * as secretsmanager from "aws-cdk-lib/aws-secretsmanager";
import * as elbv2 from "aws-cdk-lib/aws-elasticloadbalancingv2";
import * as codebuild from "aws-cdk-lib/aws-codebuild";
import * as codepipeline from "aws-cdk-lib/aws-codepipeline";
import * as codepipelineActions from "aws-cdk-lib/aws-codepipeline-actions";
import * as route53 from "aws-cdk-lib/aws-route53";
import * as targets from "aws-cdk-lib/aws-route53-targets";
import * as acm from "aws-cdk-lib/aws-certificatemanager";
import * as rds from "aws-cdk-lib/aws-rds";
import { ManagedPolicy, Role, ServicePrincipal } from "aws-cdk-lib/aws-iam";

export class AwsStack extends cdk.Stack {
    constructor(scope: Construct, id: string, props?: cdk.StackProps) {
        super(scope, id, props);

        // Misc
        const hostedZone = route53.HostedZone.fromHostedZoneAttributes(this, "appDomain", { zoneName: this.node.getContext("hostedZoneName"), hostedZoneId: this.node.getContext("hostedZoneId") });

        const ecrRepo = new ecr.Repository(this, "appEcrRepo", {
            repositoryName: "glimpz",
        });

        // Network
        const vpc = new ec2.Vpc(this, "appVpc", {
            ipAddresses: ec2.IpAddresses.cidr("10.0.0.0/16"),
            natGateways: 1,
        });

        const secret = new secretsmanager.Secret(this, "appSecret");

        const taskExecRole = new Role(this, "appTaskExecRole", {
            assumedBy: new ServicePrincipal("ecs-tasks.amazonaws.com"),
        });

        taskExecRole.addManagedPolicy(ManagedPolicy.fromAwsManagedPolicyName("service-role/AmazonECSTaskExecutionRolePolicy"));
        secret.grantRead(taskExecRole);

        // Database
        const PORT = 5432;

        const dbCluster = new rds.DatabaseCluster(this, "appDbCluster", {
            engine: rds.DatabaseClusterEngine.auroraPostgres({
                version: rds.AuroraPostgresEngineVersion.VER_14_4,
            }),
            instances: 1,
            instanceProps: {
                vpc,
                instanceType: new ec2.InstanceType("serverless"),
            },
            port: PORT,
        });

        cdk.Aspects.of(dbCluster).add({
            visit(node) {
                if (node instanceof rds.CfnDBCluster) {
                    node.serverlessV2ScalingConfiguration = {
                        minCapacity: 0.5,
                        maxCapacity: 1,
                    };
                }
            },
        });

        dbCluster.secret!.grantRead(taskExecRole);

        const bastion = new ec2.BastionHostLinux(this, "appBastionHost", {
            vpc,
        });

        dbCluster.connections.allowDefaultPortFrom(bastion);

        // Setup ECS cluster
        const taskDefinition = new ecs.FargateTaskDefinition(this, "appTaskDefinition", {
            cpu: 256,
            memoryLimitMiB: 512,
            executionRole: taskExecRole,
        });

        taskDefinition.addContainer("appContainer", {
            image: ecs.ContainerImage.fromEcrRepository(ecrRepo, "latest"),
            portMappings: [{ containerPort: 8080 }],
            environment: {
                AWS_SECRET_NAME: secret.secretName,
                DB_SECRET_NAME: dbCluster.secret!.secretName,
            },
            logging: ecs.LogDrivers.awsLogs({
                streamPrefix: "app",
            }),
        });

        const cluster = new ecs.Cluster(this, "appCluster", {
            vpc,
        });

        const fargateService = new ecs.FargateService(this, "appFargateService", {
            cluster,
            taskDefinition,
            capacityProviderStrategies: [
                {
                    capacityProvider: "FARGATE_SPOT",
                    weight: 1,
                },
            ],
            desiredCount: this.node.getContext("desiredCount"),
        });

        dbCluster.connections.allowDefaultPortFrom(fargateService);

        const loadBalancer = new elbv2.ApplicationLoadBalancer(this, "appLoadBalancer", {
            vpc,
            internetFacing: true,
        });

        new route53.ARecord(this, "appLbSubdomainRecord", {
            zone: hostedZone,
            target: route53.RecordTarget.fromAlias(new targets.LoadBalancerTarget(loadBalancer)),
            recordName: this.node.getContext("lbSubdomainName"),
        });

        const certificate = new acm.DnsValidatedCertificate(this, "myAppCertificate", {
            domainName: this.node.getContext("lbSubdomainName"),
            hostedZone: hostedZone,
        });

        const listener = loadBalancer.addListener("appListener", {
            port: 443,
            open: true,
            certificates: [certificate],
        });

        listener.addTargets("appListenerTargetGroup", {
            port: 80,
            targets: [
                fargateService.loadBalancerTarget({
                    containerName: "appContainer",
                    containerPort: 8080,
                }),
            ],
        });

        // CICD pipeline
        const buildProject = new codebuild.PipelineProject(this, "appBuildProject", {
            environment: {
                buildImage: codebuild.LinuxBuildImage.STANDARD_1_0,
                privileged: true,
            },
            environmentVariables: {
                ECR_REPO_URI: {
                    value: ecrRepo.repositoryUri,
                    type: codebuild.BuildEnvironmentVariableType.PLAINTEXT,
                },
            },
        });

        ecrRepo.grantPullPush(buildProject);

        const sourceOutput = new codepipeline.Artifact();
        const sourceAction = new codepipelineActions.GitHubSourceAction({
            actionName: "GitHubSource",
            owner: this.node.getContext("ghOwner"),
            repo: this.node.getContext("ghRepo"),
            branch: this.node.getContext("ghBranch"),
            output: sourceOutput,
            oauthToken: cdk.SecretValue.secretsManager(this.node.getContext("ghTokenSecret")),
        });

        const buildOutput = new codepipeline.Artifact();
        const buildAction = new codepipelineActions.CodeBuildAction({
            actionName: "DockerBuild",
            project: buildProject,
            input: sourceOutput,
            outputs: [buildOutput],
        });

        const deploymentAction = new codepipelineActions.EcsDeployAction({
            actionName: "DeployAction",
            service: fargateService,
            input: buildOutput,
        });

        new codepipeline.Pipeline(this, "appPipeline", {
            stages: [
                { stageName: "Source", actions: [sourceAction] },
                { stageName: "Build", actions: [buildAction] },
                { stageName: "Deploy", actions: [deploymentAction] },
            ],
        });
    }
}
