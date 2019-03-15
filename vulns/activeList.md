# pm2077 Active Vulnerabilities


XSS Stored

1.	System Notifications
    1.	Title: Executes in the header notification menu
2.	Project Task
    1.	Dashboard: Executes in the Your Open Task List

XSS Reflective

1.	Login
    1.	Query string “err” parameter
    2.	Password reset “err” parameter

CSRF

1.	Admin System Notifications
2.	PM Project Task
3.	Should be on all application forms

IDOR

1.	PM Project Task
    1.	On another user’s task

Missing Function Level Access Control

1.	PM Project Task
    1.	No PM role validation
2.	Admin Settings
    1.	No Admin role validation
3.	Customer Delete
    1.	There is no authentication checks on route “/customer/delete/{id:[0-9]+}”

User Enumeration

1.	Login
2. Forgot password page.

Sensitive Data Exposure

1.	/public/
    1.	Directory listing enabled
2.	/resources/
    1.	No authentication checks
3.	Admin Application User
    1.	Password is visible by viewing source

Command Injection

1.	Admin System Backup

Unvalidated Redirect

1.	Dashboard, Internal Resource List
2.	Admin, Internal Resources
3.	Login page, “redirect” parameter.

Path Transversal

1.	Project Task File Upload

SQLi

1.	Not working yet.  :(

Weak Password Policy
1.	Change password 6 chars 1alpha 1 number

