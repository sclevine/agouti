package appium_test

import (
        "testing"

        . "github.com/onsi/ginkgo"
        . "github.com/onsi/gomega"
)

func TestAppium(t *testing.T) {
        RegisterFailHandler(Fail)
        RunSpecs(t, "Appium Suite")
}
