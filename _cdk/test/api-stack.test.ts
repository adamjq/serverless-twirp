import { App } from "aws-cdk-lib";
import { ApiStack } from "../lib/api-stack";

describe("ApiStack", () => {
  it("synthesizes", () => {
    const app = new App({ outdir: "./cdk.out" });

    new ApiStack(app, "ApiStack", {
      env: {
        region: "us-east-1",
        account: "111111111111",
      },
    });

    app.synth();
  });
});
