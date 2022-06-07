import { Bucket, StackContext, Table } from "@serverless-stack/resources";
import { RemovalPolicy } from "aws-cdk-lib";

export function StorageStack({ stack }: StackContext) {
  // Create the table
  const table = new Table(stack, "Main", {
    fields: {
      PK: "string",
      SK: "string",
      GSI1_PK: "string",
      GS1_SK: "string",
    },
    primaryIndex: { partitionKey: "PK", sortKey: "SK" },
    globalIndexes: {
      GSI1: { partitionKey: "GSI1_PK", sortKey: "GS1_SK" },
    },
    cdk: {
      table: {
        removalPolicy: RemovalPolicy.DESTROY,
      },
    },
  });

  const bucket = new Bucket(stack, "Uploads");

  return {
    table,
    bucket,
  };
}
