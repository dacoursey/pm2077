# pm2077 Active Vulnerabilities


XSS Stored

1.	System Notifications
    -	Title: Executes in the header notification menu
        - As admin, enter a new system notification with XSS in the title:  `<script>confirm('XSS?')</script>`
2.	Project Task
    -	Dashboard: Executes in the Your Open Task List
        - As PM, enter a new task for any user with XSS in the task info: `<script>confirm('XSS?')</script>`

XSS Reflective

1.	Login
    -	Login URL string “err” parameter
        - Modify the url to include XSS in the parameter: `http://localhost:5000/login?err=<script>confirm('XSS?')</script>`
    -	Password reset “err” parameter
        - Modify the url to include XSS in the parameter: `http://localhost:5000/login?err=<script>confirm('XSS?')</script>`
2. Where else can you find "err"?

CSRF

1.	Admin System Notifications
    - POST to `/notification/save`
2.	PM Project Task
    - POST to `/task/save`
3. Should be on other forms.

IDOR

1.	Tasks
    -	View another user’s tasks
        - Logged in as a regular user, view your task list. Modify the URL bar to include a different user's ID:  `http://localhost:5000/mytask/4`
2. Look around for other usage of ID-based lookups.

Missing Function Level Access Control

1.	Tasks
    - View all tasks as regular user.
        - Logged in as a regular user, view your task list. Modify the URL bar to include a user with higher permissions ID:  `http://localhost:5000/mytask/1`
2.	Admin Settings
    - No Admin role validation on settings route.
        - Logged in as a regular user, modify the URL bar to directly browse to the admin settings: `http://localhost:5000/settings`
3.	Customer Delete
    - No Admin role validation on the delete route.
        - Logged in as a regular user, modify the URL bar to directly browse to the delete route: `http://localhost:5000/customer/delete/{id:[0-9]+}`

User Enumeration

1.	Login
    - Different error messages for login failures.
        - Log in with valid user and bad password and save the result. 
        - Log in with bad user and compare the result.

Sensitive Data Exposure

1.	/public/
    - Directory listing enabled
        - Browse directly to `http://localhost:5000/public/`.
2.	/resources/
    - No authentication checks
        - Log out of the application and clear all cookies to ensure there is  no current session.
        - Browse directly to `http://localhost:5000/resources` and view all files that have been uploaded.
3.	Passwords not protected.
    - Password is visible by viewing source
        - Browse to the user list and select "edit" to view one user. Right-click the password box and select "Inspect Element" and validate the user's password is displayed.

Command Injection

1.	Admin System Backup
    - Backup type is injectable.
        - Use burp repeater and modify the `backupType` parameter to be an OS command:   `backupType=ls`

Unvalidated Redirect

1.	Login page
    - Parameter used to redirect to the login page.
        - Modify the URL to redirect the user after login: `http://localhost:5000/login?redirect=http://www.appsec.lol`
2.	Internal Resources
    - Resource included URL parameter.

Path Transversal

1.	Project Task File Upload

SQLi

1.	Not working yet. :(

Weak Password Policy
1. Current policy: 6 chars - 1 alpha 1 number

