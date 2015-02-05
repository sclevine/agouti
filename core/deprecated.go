package core

// Chrome returns an instance of a ChromeDriver WebDriver.
// DEPRECATED: Use ChromeDriver() instead.
func Chrome() (WebDriver, error) {
	return ChromeDriver(), nil
}
