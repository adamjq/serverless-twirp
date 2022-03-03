import { Construct } from "constructs";
import { Duration, Stack, StackProps } from "aws-cdk-lib";
import { Code, Function, Runtime } from "aws-cdk-lib/aws-lambda";
import { LambdaRestApi } from "aws-cdk-lib/aws-apigateway";
import { AttributeType, BillingMode, Table } from "aws-cdk-lib/aws-dynamodb";

const LAMBDA_RUNTIME = Runtime.GO_1_X;
const LAMBDA_DEFAULT_TIMEOUT = Duration.seconds(30);
const LAMBDA_DEFAULT_MEMORY = 256; // megabytes
const CODE_DIST = "../dist";

export class ApiStack extends Stack {
  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);

    // follow single-table design pattern
    const ddbTable = new Table(this, "DynamoTable", {
      billingMode: BillingMode.PAY_PER_REQUEST,
      partitionKey: {
        name: "PK",
        type: AttributeType.STRING,
      },
      sortKey: {
        name: "SK",
        type: AttributeType.STRING,
      },
    });

    const backend = new Function(this, "ApiLambdaBackend", {
      runtime: LAMBDA_RUNTIME,
      handler: "api",
      code: Code.fromAsset(CODE_DIST),
      timeout: LAMBDA_DEFAULT_TIMEOUT,
      memorySize: LAMBDA_DEFAULT_MEMORY,
      environment: {
        BACKEND_TABLE: ddbTable.tableName,
      },
    });

    // proxy requests to lambda by default
    new LambdaRestApi(this, "Api", {
      handler: backend,
    });

    // IAM
    ddbTable.grantFullAccess(backend);
  }
}
