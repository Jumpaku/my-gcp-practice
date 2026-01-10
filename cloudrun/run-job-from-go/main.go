package main

import (
	"context"
	"fmt"
	"os"

	run "cloud.google.com/go/run/apiv2"

	"cloud.google.com/go/run/apiv2/runpb"
)

func main() {
	var projectID, region, job = os.Getenv("PROJECT_ID"), os.Getenv("REGION"), os.Getenv("JOB")
	if projectID == "" || region == "" || job == "" {
		panic("PROJECT_ID, REGION and JOB environment variables must be set.")
	}

	ctx := context.Background()
	c, err := run.NewJobsClient(ctx)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	req := &runpb.RunJobRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/jobs/%s", projectID, region, job),
		Overrides: &runpb.RunJobRequest_Overrides{
			ContainerOverrides: []*runpb.RunJobRequest_Overrides_ContainerOverride{
				{
					Args: []string{"abc", "def", "ghi"},
					Env: []*runpb.EnvVar{
						{
							Name:   "OVERRIDE_ENV",
							Values: &runpb.EnvVar_Value{Value: "jkl"},
						},
					},
				},
			},
		},
	}
	op, err := c.RunJob(ctx, req)
	if err != nil {
		panic(err)
	}

	resp, err := op.Wait(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v", resp)
}
