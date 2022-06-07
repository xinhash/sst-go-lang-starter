import { StackContext, Api, use, Auth } from "@serverless-stack/resources";
import { StorageStack } from "./StorageStack";

export function ApiStack({ stack }: StackContext) {
  const { table } = use(StorageStack);

  const auth = new Auth(stack, "Auth", {
    login: ["email"],
    cdk: {
      userPoolClient: {
        authFlows: {
          userPassword: true,
        },
      },
    },
  });

  const api = new Api(stack, "api", {
    defaults: {
      authorizer: "iam",
      function: {
        environment: {
          table: table.tableName,
          cognitoClientId: auth.userPoolClientId,
          cognitoUserPoolId: auth.userPoolId,
        },
        permissions: [table],
      },
    },
    routes: {
      "GET /resources": "functions/resources-find-all/lambda.go",
      "GET /resources/{id}": "functions/resources-find-one/lambda.go",
      "POST /signup": {
        function: "functions/auth-signup/lambda.go",
        authorizer: "none",
      },
      "POST /signin": {
        function: "functions/auth-signin/lambda.go",
        authorizer: "none",
      },
      "POST /signup/confirm": {
        function: "functions/auth-confirm/lambda.go",
        authorizer: "none",
      },
      "POST /forgot-password": {
        function: "functions/auth-forgot-password/lambda.go",
        authorizer: "none",
      },
      "POST /forgot-password/confirm": {
        function: "functions/auth-confirm-forgot-password/lambda.go",
        authorizer: "none",
      },
      "GET /": {
        function: "functions/lambda.go",
        authorizer: "none",
      },
    },
  });

  auth.attachPermissionsForAuthUsers([api, table]);

  stack.addOutputs({
    ApiEndpoint: api.url,
    UserPoolId: auth.userPoolId,
    IdentityPoolId: auth.cognitoIdentityPoolId || "",
    UserPoolClientId: auth.userPoolClientId,
  });
}
