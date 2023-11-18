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
import { ManagedPolicy, Role, ServicePrincipal } from "aws-cdk-lib/aws-iam";

export class InfraStack extends cdk.NestedStack {
    constructor(scope: Construct, id: string, props?: cdk.StackProps) {
        super(scope, id, props);

        // Define runtime infrastructure
        const ecrRepo = new ecr.Repository(this, "appEcrRepo", {
            repositoryName: "glimpz",
        });

        const vpc = new ec2.Vpc(this, "appVpc", {
            ipAddresses: ec2.IpAddresses.cidr("10.0.0.0/16"),
        });

        const secret = new secretsmanager.Secret(this, "appSecret");

        const taskExecRole = new Role(this, "appTaskExecRole", {
            assumedBy: new ServicePrincipal("ecs-tasks.amazonaws.com"),
        });

        taskExecRole.addManagedPolicy(ManagedPolicy.fromAwsManagedPolicyName("service-role/AmazonECSTaskExecutionRolePolicy"));
        secret.grantRead(taskExecRole);

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
            },
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
            desiredCount: 1,
        });

        const loadBalancer = new elbv2.ApplicationLoadBalancer(this, "appLoadBalancer", {
            vpc,
            internetFacing: true,
        });

        const listener = loadBalancer.addListener("appListener", {
            port: 80,
            open: true,
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
            owner: this.node.getContext("githubOwner"),
            repo: this.node.getContext("githubRepo"),
            branch: this.node.getContext("githubBranch"),
            output: sourceOutput,
            oauthToken: cdk.SecretValue.secretsManager("github-token"),
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

        const pipeline = new codepipeline.Pipeline(this, "appPipeline", {
            stages: [
                { stageName: "Source", actions: [sourceAction] },
                { stageName: "Build", actions: [buildAction] },
                { stageName: "Deploy", actions: [deploymentAction] },
            ],
        });
    }
}
