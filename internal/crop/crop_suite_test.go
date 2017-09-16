package crop_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCrop(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Crop Suite")
}
