import * as cdk from "aws-cdk-lib";
import { Construct } from "constructs";
import { InfraStack } from "./infra-stack";

export class AwsStack extends cdk.Stack {
    constructor(scope: Construct, id: string, props?: cdk.StackProps) {
        super(scope, id, props);

        new InfraStack(this, "infraStack", {});
    }
}
