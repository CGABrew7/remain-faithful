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
        // Dismiss any iOS system sheet (Strong Password suggestion, AutoFill, etc.)
        addUIInterruptionMonitor(withDescription: "System sheet") { sheet in
            for label in ["Cancel", "Not Now", "Use Strong Password", "Choose My Own Password"] {
                if sheet.buttons[label].exists {
                    sheet.buttons[label].tap()
                    return true
                }
            }
            return false
        }

        app.buttons["Begin Your Journey"].tap()

        let getStarted = app.buttons["Get Started"]
        XCTAssertTrue(getStarted.waitForExistence(timeout: 3))
        XCTAssertFalse(getStarted.isEnabled, "button should be disabled with empty fields")

        let nameField = app.textFields.matching(identifier: "name-field").firstMatch
        XCTAssertTrue(nameField.exists, "name-field must exist")
        nameField.tap()
        nameField.typeText("Alice")

        let emailField = app.textFields.matching(identifier: "email-field").firstMatch
        emailField.tap()
        emailField.typeText("a@b.co")

        let pwField = app.secureTextFields.matching(identifier: "password-field").firstMatch
        XCTAssertTrue(pwField.exists, "password-field must exist")
        pwField.tap()
        pwField.typeText("Passw0rd")  // 8 chars

        let enabledPredicate = NSPredicate(format: "isEnabled == YES")
        let result = XCTWaiter().wait(
            for: [expectation(for: enabledPredicate, evaluatedWith: getStarted)],
            timeout: 3
        )
        XCTAssertEqual(result, .completed, "button should be enabled with valid name, email, and 8+ char password")
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
