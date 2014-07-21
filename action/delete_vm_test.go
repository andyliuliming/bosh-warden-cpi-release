package action_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "bosh-warden-cpi/action"
	fakevm "bosh-warden-cpi/vm/fakes"
)

var _ = Describe("DeleteVM", func() {
	var (
		vmFinder *fakevm.FakeFinder
		action   DeleteVM
	)

	BeforeEach(func() {
		vmFinder = &fakevm.FakeFinder{}
		action = NewDeleteVM(vmFinder)
	})

	Describe("Run", func() {
		It("tries to find vm with given vm cid", func() {
			_, err := action.Run("fake-vm-id")
			Expect(err).ToNot(HaveOccurred())

			Expect(vmFinder.FindID).To(Equal("fake-vm-id"))
		})

		Context("when vm is found with given vm cid", func() {
			var (
				vm *fakevm.FakeVM
			)

			BeforeEach(func() {
				vm = fakevm.NewFakeVM("fake-vm-id")
				vmFinder.FindVM = vm
				vmFinder.FindFound = true
			})

			It("deletes vm", func() {
				_, err := action.Run("fake-vm-id")
				Expect(err).ToNot(HaveOccurred())

				Expect(vm.DeleteCalled).To(BeTrue())
			})

			It("returns error if deleting vm fails", func() {
				vm.DeleteErr = errors.New("fake-delete-err")

				_, err := action.Run("fake-vm-id")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fake-delete-err"))
			})
		})

		Context("when vm is not found with given cid", func() {
			It("does vmFinder return error", func() {
				vmFinder.FindFound = false

				_, err := action.Run("fake-vm-id")
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when vm finding fails", func() {
			It("does not return error", func() {
				vmFinder.FindErr = errors.New("fake-find-err")

				_, err := action.Run("fake-vm-id")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fake-find-err"))
			})
		})
	})
})
