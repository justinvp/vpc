import * as vpc from "@pulumi/vpc";

const res = new vpc.Vpc("my-vpc", {
    tags: {
        foo: "bar",
    },
});
