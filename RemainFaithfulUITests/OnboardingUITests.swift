import XCTest

final class OnboardingUITests: XCTestCase {

    var app: XCUIApplication!

    override func setUp() {
        super.setUp()
        continueAfterFailure = false
        app = XCUIApplication()
        app.launchArguments = ["-UITestReset"]
        app.launch()
    }

    override func tearDown() {
        app.terminate()
        app = nil
        super.tearDown()
    }

    // MARK: - Test 1: WelcomeStep renders correctly

    func testWelcomeStep_titleAndCTAVisible() {
        XCTAssertTrue(app.staticTexts["Remain Faithful"].waitForExistence(timeout: 3))
        XCTAssertTrue(app.buttons["Begin Your Journey"].exists)
        XCTAssertTrue(app.buttons["Begin Your Journey"].isEnabled)
    }

    // MARK: - Test 2: Registration validation gate

    func testCreateAccountStep_getStartedDisabledUntilValidInput() {
        app.buttons["Begin Your Journey"].tap()

        let getStarted = app.buttons["Get Started"]
        XCTAssertTrue(getStarted.waitForExistence(timeout: 3))
        XCTAssertFalse(getStarted.isEnabled, "button should be disabled with empty fields")

        // Fill name
        let nameField = app.textFields.element(boundBy: 0)
        nameField.tap()
        nameField.typeText("Alice")

        // Fill valid email
        let emailField = app.textFields.element(boundBy: 1)
        emailField.tap()
        emailField.typeText("alice@example.com")

        // Password too short — still disabled
        let passwordField = app.secureTextFields.firstMatch
        passwordField.tap()
        passwordField.typeText("short")
        XCTAssertFalse(getStarted.isEnabled, "button should remain disabled with < 8 char password")

        // Replace short password with a valid one
        passwordField.tap()
        passwordField.typeText(String(repeating: XCUIKeyboardKey.delete.rawValue, count: 5))
        passwordField.typeText("Password1!")
        XCTAssertTrue(getStarted.isEnabled, "button should be enabled with valid name, email, and 8+ char password")
    }

    // MARK: - Test 3: Sign In link shows LoginView

    func testSignInLink_presentsLoginView() {
        // Tap the "Sign In" portion of the WelcomeStep footer button
        let signInButton = app.buttons.matching(
            NSPredicate(format: "label CONTAINS[c] 'Sign In'")
        ).firstMatch
        XCTAssertTrue(signInButton.waitForExistence(timeout: 3))
        signInButton.tap()

        // LoginView shows "Welcome Back" heading and an email text field
        XCTAssertTrue(
            app.staticTexts["Welcome Back"].waitForExistence(timeout: 3),
            "LoginView should show 'Welcome Back' heading"
        )
        XCTAssertTrue(app.textFields.firstMatch.exists)
    }
}
