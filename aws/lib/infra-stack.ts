import * as cdk from "aws-cdk-lib";
import { Construct } from "constructs";
import * as ecr from "aws-cdk-lib/aws-ecr";
import * as ecs from "aws-cdk-lib/aws-ecs";
import * as ec2 from "aws-cdk-lib/aws-ec2";
import * as secretsmanager from "aws-cdk-lib/aws-secretsmanager";

export class InfraStack extends cdk.NestedStack {
    constructor(scope: Construct, id: string, props?: cdk.StackProps) {
        super(scope, id, props);

        const ecrRepo = ecr.Repository.fromRepositoryArn(this, "appEcrRepo", this.node.getContext("ecrRepo"));

        const vpc = new ec2.Vpc(this, "appVpc", {
            ipAddresses: ec2.IpAddresses.cidr("10.0.0.0/16"),
        });

        const secret = new secretsmanager.Secret(this, "appSecret");

        const taskDefinition = new ecs.FargateTaskDefinition(this, "appTaskDefinition", {
            cpu: 256,
            memoryLimitMiB: 512,
        });

        taskDefinition.addContainer("appContainer", {
            image: ecs.ContainerImage.fromEcrRepository(ecrRepo, "latest"),
            portMappings: [{ containerPort: 80 }],
            environment: {},
        });

        const cluster = new ecs.Cluster(this, "appCluster", {
            vpc,
        });

        // **** I need to assign permissions here for my container to access secret variables

        const fargateService = new ecs.FargateService(this, "appFargateService", {
            cluster,
            taskDefinition,
            capacityProviderStrategies: [
                {
                    capacityProvider: "FARGATE_SPOT",
                    weight: 1,
                },
            ],
        });
    }
}
