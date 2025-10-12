0. Read this guide before https://github.com/who-knows-inc/EK_DAT_DevOps_2025_Autumn/blob/main/00._Course_Material/01._Assignments/05._Docker_The_Simulation/02._After/setup_postman_monitoring.md


1. Create Postman Environment

Go to Environments → Manage Environments → Add:

Variable	Initial Value	        Notes
IP	        68.221.201.252	        API host
PORT	    8080	                API port

Call enviroment API Tests

2. Create Collection & Requests

Create a collection API Tests. F.eks register example

POST /api/register

Body:
{
  "username": "{{username}}",
  "email": "{{email}}",
  "password": "mySecret123",
  "password2": "mySecret123"
}


Pre-request Script:
const randomId = Math.floor(Math.random() * 1000000);
pm.environment.set("username", `testuser_${randomId}`);
pm.environment.set("email", `testuser_${randomId}@example.com`);

Tests:
pm.test("Registration succeeded", () => {
    pm.response.to.have.status(201);
    const json = pm.response.json();
    pm.expect(json.username).to.eql(pm.environment.get("username"));
});


4. Run with Postman Runner

Open Collection Runner → select API Tests

Select environment API Test

Optionally, run multiple iterations to test different users

Click Run

Setup scheduled run