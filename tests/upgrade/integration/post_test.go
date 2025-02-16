// go:build upgrade_integration

package integration

import (
	"fmt"
	"path/filepath"

	clihelpers "github.com/cheqd/cheqd-node/tests/integration/helpers"
	cli "github.com/cheqd/cheqd-node/tests/upgrade/integration/cli"
	didcli "github.com/cheqd/cheqd-node/x/did/client/cli"
	didtypesv2 "github.com/cheqd/cheqd-node/x/did/types"
	resourcetypesv2 "github.com/cheqd/cheqd-node/x/resource/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Upgrade - Post", func() {
	var feeParams didtypesv2.FeeParams
	var resourceFeeParams resourcetypesv2.FeeParams

	BeforeEach(func() {
		// Query fee params
		res, err := cli.QueryParams(cli.Validator0, didtypesv2.ModuleName, string(didtypesv2.ParamStoreKeyFeeParams))
		Expect(err).To(BeNil())
		err = clihelpers.Codec.UnmarshalJSON([]byte(res.Value), &feeParams)
		Expect(err).To(BeNil())

		res, err = cli.QueryParams(cli.Validator0, resourcetypesv2.ModuleName, string(resourcetypesv2.ParamStoreKeyFeeParams))
		Expect(err).To(BeNil())
		err = clihelpers.Codec.UnmarshalJSON([]byte(res.Value), &resourceFeeParams)
		Expect(err).To(BeNil())
	})

	Context("After a software upgrade execution has concluded", func() {
		It("should wait for node catching up", func() {
			By("pinging the node status until catching up is flagged as false")
			err := cli.WaitForCaughtUp(cli.Validator0, cli.CliBinaryName, cli.VotingPeriod*6)
			Expect(err).To(BeNil())
		})

		It("should wait for a certain number of blocks to be produced", func() {
			By("fetching the current chain height")
			currentHeight, err := cli.GetCurrentBlockHeight(cli.Validator0, cli.CliBinaryName)
			Expect(err).To(BeNil())

			By("waiting for 10 blocks to be produced on top, after the upgrade")
			err = cli.WaitForChainHeight(cli.Validator0, cli.CliBinaryName, currentHeight+10, cli.VotingPeriod*6)
			Expect(err).To(BeNil())
		})

		It("should match the expected module version map", func() {
			By("loading the expected module version map")
			var expected upgradetypes.QueryModuleVersionsResponse
			_, err := Loader(filepath.Join(GeneratedJSONDir, "post", "query - module-version-map", "v1.json"), &expected)
			Expect(err).To(BeNil())

			By("matching the expected module version map")
			actual, err := cli.QueryModuleVersionMap(cli.Validator0)
			Expect(err).To(BeNil())

			Expect(actual.ModuleVersions).To(Equal(expected.ModuleVersions), "module version map mismatch")
		})

		It("should load and run existing diddoc payloads - case: update", func() {
			By("matching the glob pattern for existing diddoc payloads")
			DidDocUpdatePayloads, err := RelGlob(GeneratedJSONDir, "post", "update - diddoc", "*.json")
			Expect(err).To(BeNil())

			for _, payload := range DidDocUpdatePayloads {
				var DidDocUpdatePayload didcli.DIDDocument
				var DidDocUpdateSignInput []didcli.SignInput

				testCase := GetCaseName(payload)
				By("Running: " + testCase)
				fmt.Println("Running: " + testCase)

				By("reading ")
				DidDocUpdateSignInput, err = Loader(payload, &DidDocUpdatePayload)
				Expect(err).To(BeNil())

				tax := feeParams.UpdateDid.String()
				res, err := cli.UpdateDid(DidDocUpdatePayload, DidDocUpdateSignInput, cli.Validator0, "", tax)
				Expect(err).To(BeNil())
				Expect(res.Code).To(BeEquivalentTo(0))
			}
		})

		It("should load and run existing diddoc payloads - case: deactivate", func() {
			By("matching the glob pattern for existing diddoc payloads")
			DidDocDeactivatePayloads, err := RelGlob(GeneratedJSONDir, "post", "deactivate - diddoc", "*.json")
			Expect(err).To(BeNil())

			for _, payload := range DidDocDeactivatePayloads {
				var DidDocDeacctivatePayload didtypesv2.MsgDeactivateDidDocPayload
				var DidDocDeactivateSignInput []didcli.SignInput

				testCase := GetCaseName(payload)
				By("Running: " + testCase)
				fmt.Println("Running: " + testCase)

				By("reading ")
				DidDocDeactivateSignInput, err = Loader(payload, &DidDocDeacctivatePayload)
				Expect(err).To(BeNil())

				tax := feeParams.DeactivateDid.String()
				res, err := cli.DeactivateDid(DidDocDeacctivatePayload, DidDocDeactivateSignInput, cli.Validator0, tax)
				Expect(err).To(BeNil())
				Expect(res.Code).To(BeEquivalentTo(0))
			}
		})

		It("should create resources after upgrade", func() {
			By("matching the glob pattern for resource payloads to create")
			ResourcePayloads, err := RelGlob(GeneratedJSONDir, "post", "create - resource", "*.json")
			Expect(err).To(BeNil())

			for _, payload := range ResourcePayloads {
				var ResourceCreatePayload resourcetypesv2.MsgCreateResourcePayload

				testCase := GetCaseName(payload)
				By("Running: create " + testCase)
				fmt.Println("Running: " + testCase)

				signInputs, err := Loader(payload, &ResourceCreatePayload)
				Expect(err).To(BeNil())

				ResourceFile, err := CreateTestJSON(GinkgoT().TempDir(), ResourceCreatePayload.Data)
				Expect(err).To(BeNil())

				res, err := cli.CreateResource(
					ResourceCreatePayload,
					ResourceFile,
					signInputs,
					cli.Validator0,
					resourceFeeParams.Json.String(),
				)

				Expect(err).To(BeNil())
				Expect(res.Code).To(Equal(uint32(0)))
			}
		})

		It("should load and run expected diddoc payloads", func() {
			By("matching the glob pattern for existing diddoc payloads")
			ExpectedDidDocUpdateRecords, err := RelGlob(GeneratedJSONDir, "post", "query - diddoc", "*.json")
			Expect(err).To(BeNil())

			for _, payload := range ExpectedDidDocUpdateRecords {
				var DidDocUpdateRecord didtypesv2.DidDoc

				testCase := GetCaseName(payload)
				By("Running: query " + testCase)
				fmt.Println("Running: " + testCase)

				_, err = Loader(payload, &DidDocUpdateRecord)
				Expect(err).To(BeNil())

				res, err := cli.QueryDid(DidDocUpdateRecord.Id, cli.Validator0)
				Expect(err).To(BeNil())

				if DidDocUpdateRecord.Context == nil {
					DidDocUpdateRecord.Context = []string{}
				}
				if DidDocUpdateRecord.Authentication == nil {
					DidDocUpdateRecord.Authentication = []string{}
				}
				if DidDocUpdateRecord.AssertionMethod == nil {
					DidDocUpdateRecord.AssertionMethod = []string{}
				}
				if DidDocUpdateRecord.CapabilityInvocation == nil {
					DidDocUpdateRecord.CapabilityInvocation = []string{}
				}
				if DidDocUpdateRecord.CapabilityDelegation == nil {
					DidDocUpdateRecord.CapabilityDelegation = []string{}
				}
				if DidDocUpdateRecord.KeyAgreement == nil {
					DidDocUpdateRecord.KeyAgreement = []string{}
				}
				if DidDocUpdateRecord.Service == nil {
					DidDocUpdateRecord.Service = []*didtypesv2.Service{}
				}
				if DidDocUpdateRecord.AlsoKnownAs == nil {
					DidDocUpdateRecord.AlsoKnownAs = []string{}
				}

				Expect(*res.Value.DidDoc).To(Equal(DidDocUpdateRecord))
			}
		})

		It("should load and run expected resource payloads", func() {
			By("matching the glob pattern for existing resource payloads")
			ExpectedResourceCreateRecords, err := RelGlob(GeneratedJSONDir, "post", "query - resource", "*.json")
			Expect(err).To(BeNil())

			for _, payload := range ExpectedResourceCreateRecords {
				var ResourceCreateRecord resourcetypesv2.ResourceWithMetadata

				testCase := GetCaseName(payload)
				By("Running: query " + testCase)
				fmt.Println("Running: " + testCase)

				_, err = Loader(payload, &ResourceCreateRecord)
				Expect(err).To(BeNil())

				res, err := cli.QueryResource(ResourceCreateRecord.Metadata.CollectionId, ResourceCreateRecord.Metadata.Id, cli.Validator0)

				Expect(err).To(BeNil())
				Expect(res.Resource.Metadata.Id).To(Equal(ResourceCreateRecord.Metadata.Id))
				Expect(res.Resource.Metadata.CollectionId).To(Equal(ResourceCreateRecord.Metadata.CollectionId))
				Expect(res.Resource.Metadata.Name).To(Equal(ResourceCreateRecord.Metadata.Name))
				Expect(res.Resource.Metadata.Version).To(Equal(ResourceCreateRecord.Metadata.Version))
				Expect(res.Resource.Metadata.ResourceType).To(Equal(ResourceCreateRecord.Metadata.ResourceType))
				Expect(res.Resource.Metadata.AlsoKnownAs).To(Equal(ResourceCreateRecord.Metadata.AlsoKnownAs))
				Expect(res.Resource.Metadata.MediaType).To(Equal(ResourceCreateRecord.Metadata.MediaType))
				// Created is populated on successful creation. We are ignoring it here.
				// Expect(res.Resource.Metadata.Created).To(Equal(ResourceCreateRecord.Metadata.Created))
				Expect(res.Resource.Metadata.Checksum).To(Equal(ResourceCreateRecord.Metadata.Checksum))
				Expect(res.Resource.Metadata.PreviousVersionId).To(Equal(ResourceCreateRecord.Metadata.PreviousVersionId))
				Expect(res.Resource.Metadata.NextVersionId).To(Equal(ResourceCreateRecord.Metadata.NextVersionId))
			}
		})
	})
})
