# **Business Requirements Document (BRD)**

**Project Name:** HR Claims Management App
**Version:** 1.0
**Author:** \[Your Name]
**Date:** 2025-05-27

---

## **1. Purpose**

The purpose of this project is to develop an HR Claims Management application that enables employees to submit and track claims, and allows HR administrators to manage and process these claims through a structured approval workflow.

---

## **2. Scope**

This system will support:

* Claim submission and tracking for normal users (employees)
* Claim type management, user role/group management, and claim approval workflows for admin users
* A multi-level approval flow per user group
* Role-based access and actions

---

## **3. User Roles & Permissions**

### **Normal Users (Employees)**

* Submit claims
* Cancel their own claim submission (before final approval)
* View all their claims and statuses (including rejected)

### **Admin Users (HR/Admin)**

* Create/edit/delete (soft) **claim types**
* View all claims in the system
* Submit and cancel claims (only for themselves)
* Approve/reject claims (except their own)
* Create/update/delete (soft) **user groups**
* Change user type (normal ↔ admin)
* Create approval levels per user group
* Assign/remove approvers for each level and define which claim status transitions they can perform (e.g., approve, reject)

---

## **4. Claim Workflow**

**Statuses:**

```
Draft → Submitted → Approved → Payment-in-progress → Paid  
                      ↳ Rejected
```

* Claims are initially saved as **Draft**
* Upon submission, status changes to **Submitted**
* Based on defined approval levels, approvers can **Approve** or **Reject**
* Once fully approved, status transitions to **Payment-in-progress**, then to **Paid**

---

## **5. Functional Requirements**

| Feature              | Description                                                       | Role                            |
| -------------------- | ----------------------------------------------------------------- | ------------------------------- |
| Submit Claim         | Fill out form with claim details & submit                         | Normal, Admin (self only)       |
| Cancel Claim         | Cancel unapproved claim                                           | Normal, Admin (self only)       |
| View Claims          | List all personal claims with statuses                            | All                             |
| View All Claims      | View and filter all submitted claims                              | Admin only                      |
| Create Claim Types   | Define various claim categories (e.g., Travel, Medical)           | Admin only                      |
| User Role Management | Promote/demote users between normal and admin                     | Admin only                      |
| Group Management     | Create/update/delete user groups                                  | Admin only                      |
| Approval Workflow    | Configure approval levels per group, assign approvers and actions | Admin only                      |
| Approve/Reject       | Perform actions based on level permissions                        | Admin only (not for own claims) |

---

## **6. Tech Stack**

### **Backend**

* **Language:** Golang
* **Router:** Chi
* **ORM:** GORM
* **Database:** PostgreSQL

### **Frontend**

* **Framework:** Vue 3 (Composition API)
* **Language:** TypeScript
* **State Management:** Pinia
* **Routing:** Vue Router
* **HTTP:** Axios
* **UI Design:** Sleek, professional and classy with PrimeVue library

---

## **7. Non-Functional Requirements**

* **Responsiveness:** Fully responsive for desktop and mobile devices
* **Security:** Role-based access control
* **Data Integrity:** Claims and user data should be validated and consistently stored
* **Error Handling:** Clear and descriptive error messages across system
* **Performance:** Fast and minimal-latency interface and API responses
* **MVP Standard:** All functionalities listed must be working correctly and error-free

---

## **8. Deliverables**

* Fully functional HR Claims Management application (frontend + backend)
* Admin dashboard
* Employee dashboard
* User group and approval workflow configuration module
* Deployment-ready codebase with documentation

