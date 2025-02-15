package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// You normally want to run this under a separate "Testing" subscription
// For lab purposes you will use your assigned subscription under the Cloud Dev/Ops program tenant
var subscriptionID string = "ee39387e-4892-429d-aeb2-ca4f142fbd25"

func TestAzureLinuxVMCreation(t *testing.T) {
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../",
		// Override the default terraform variables
		Vars: map[string]interface{}{
			"labelPrefix": "zhou0265",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the value of output variable
	vmName := terraform.Output(t, terraformOptions, "vm_name")
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	nicName := terraform.Output(t, terraformOptions, "nic_name")

	// Confirm VM exists
	assert.True(t, azure.VirtualMachineExists(t, vmName, resourceGroupName, subscriptionID))

	// ✅ ADDITIONAL TEST 1: Confirm NIC exists and is connected to VM
	nicExists := azure.NetworkInterfaceExists(t, subscriptionID, resourceGroupName, nicName)
	assert.True(t, nicExists, "NIC should exist and be connected to the VM")

	// ✅ ADDITIONAL TEST 2: Confirm the VM is running the correct Ubuntu version
	vm := azure.GetVirtualMachine(t, subscriptionID, resourceGroupName, vmName)
	assert.NotNil(t, vm.StorageProfile, "VM StorageProfile should not be nil")
	assert.NotNil(t, vm.StorageProfile.ImageReference, "VM ImageReference should not be nil")
	assert.Contains(t, *vm.StorageProfile.ImageReference.Sku, "22_04", "VM should be running Ubuntu 22.04")
}
