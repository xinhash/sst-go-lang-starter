import { ApiStack } from "./ApiStack";
import { App } from "@serverless-stack/resources";
import { StorageStack } from "./StorageStack";

export default function (app: App) {
  // FIXME: This should be used with grain of salt
  // Remove all resources when the dev stage is removed
  if (app.stage === "dev") {
    app.setDefaultRemovalPolicy("destroy");
  }

  app.setDefaultFunctionProps({
    runtime: "go1.x",
    srcPath: "backend",
  });
  app.stack(StorageStack);
  app.stack(ApiStack);
}
