package aws

import (
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"
)

// ACC test that's not ready
// func TestAccAWSConfig_basic(t *testing.T) {
// 	fmt.Println("Testing AWS Config")
// 	var v ec2.Vpc
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:     func() { testAccPreCheck(t) },
// 		Providers:    testAccProviders,
// 		CheckDestroy: testAccCheckInstanceDestroy,
// 		Steps: []resource.TestStep{
// 			resource.TestStep{
// 				Config: testAccVpcConfigTest,
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckVpcExists("aws_vpc.foo", &v),
// 				),
// 			},
// 		},
// 	})
// }

// const testAccConfig_pre = `
// provider "aws" {
// 	access_key = "thing"
// 	secret_key = "thing"
// }
// `

// const testAccVpcConfigTest = `
// provider "aws" {
// 	access_key = "thing"
// 	secret_key = "thing"
// }
// resource "aws_vpc" "foo" {
// 	cidr_block = "10.1.0.0/16"
// }
// `

var k = os.Getenv("AWS_ACCESS_KEY_ID")
var s = os.Getenv("AWS_SECRET_ACCESS_KEY")
var to = os.Getenv("AWS_SESSION_TOKEN")

func TestAWSConfig_shouldError(t *testing.T) {
	unsetEnv(t)
	defer resetEnv(t)
	cfg := Config{}

	c := getCreds(cfg.AccessKey, cfg.SecretKey, cfg.Token)
	_, err := c.Get()
	if awsErr, ok := err.(awserr.Error); ok {
		if awsErr.Code() != "NoCredentialProviders" {
			t.Fatalf("Expected NoCredentialProviders error")
		}
	}
	if err == nil {
		t.Fatalf("Expected an error with empty env, keys, and IAM in AWS Config")
	}
}

func TestAWSConfig_shouldBeStatic(t *testing.T) {
	unsetEnv(t)
	defer resetEnv(t)
	simple := []struct {
		Key, Secret, Token string
	}{
		{
			Key:    "test",
			Secret: "secret",
		}, {
			Key:    "test",
			Secret: "test",
			Token:  "test",
		},
		//{},
	}

	for _, c := range simple {
		cfg := Config{
			AccessKey: c.Key,
			SecretKey: c.Secret,
			Token:     c.Token,
		}

		creds := getCreds(cfg.AccessKey, cfg.SecretKey, cfg.Token)
		if creds == nil {
			t.Fatalf("Expected a static creds provider to be returned")
		}
		v, err := creds.Get()
		if err != nil {
			t.Fatalf("Error gettings creds: %s", err)
		}
		if v.AccessKeyID != c.Key {
			t.Fatalf("AccessKeyID mismatch, expected: (%s), got (%s)", c.Key, v.AccessKeyID)
		}
		if v.SecretAccessKey != c.Secret {
			t.Fatalf("SecretAccessKey mismatch, expected: (%s), got (%s)", c.Secret, v.SecretAccessKey)
		}
		if v.SessionToken != c.Token {
			t.Fatalf("SessionToken mismatch, expected: (%s), got (%s)", c.Token, v.SessionToken)
		}
	}
}

func TestAWSConfig_shouldBeENV(t *testing.T) {
	// ENV should be set
	cfg := Config{}

	creds := getCreds(cfg.AccessKey, cfg.SecretKey, cfg.Token)
	if creds == nil {
		t.Fatalf("Expected a static creds provider to be returned")
	}
	v, err := creds.Get()
	if err != nil {
		t.Fatalf("Error gettings creds: %s", err)
	}
	if v.AccessKeyID != k {
		t.Fatalf("AccessKeyID mismatch, expected: (%s), got (%s)", k, v.AccessKeyID)
	}
	if v.SecretAccessKey != s {
		t.Fatalf("SecretAccessKey mismatch, expected: (%s), got (%s)", s, v.SecretAccessKey)
	}
	if v.SessionToken != to {
		t.Fatalf("SessionToken mismatch, expected: (%s), got (%s)", to, v.SessionToken)
	}
}

// unsetEnv unsets enviornment variables for testing a "clean slate" with no
// credentials in the environment
func unsetEnv(t *testing.T) {
	if err := os.Unsetenv("AWS_ACCESS_KEY_ID"); err != nil {
		t.Fatalf("Error unsetting env var AWS_ACCESS_KEY_ID: %s", err)
	}
	if err := os.Unsetenv("AWS_SECRET_ACCESS_KEY"); err != nil {
		t.Fatalf("Error unsetting env var AWS_SECRET_ACCESS_KEY: %s", err)
	}
	if err := os.Unsetenv("AWS_SESSION_TOKEN"); err != nil {
		t.Fatalf("Error unsetting env var AWS_SESSION_TOKEN: %s", err)
	}
}

func resetEnv(t *testing.T) {
	// re-set all the envs we unset above
	if err := os.Setenv("AWS_ACCESS_KEY_ID", k); err != nil {
		t.Fatalf("Error resetting env var AWS_ACCESS_KEY_ID: %s", err)
	}
	if err := os.Setenv("AWS_SECRET_ACCESS_KEY", s); err != nil {
		t.Fatalf("Error resetting env var AWS_SECRET_ACCESS_KEY: %s", err)
	}
	if err := os.Setenv("AWS_SESSION_TOKEN", to); err != nil {
		t.Fatalf("Error resetting env var AWS_SESSION_TOKEN: %s", err)
	}
}
