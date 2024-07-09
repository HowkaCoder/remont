# remont

+---------------------------+
|          users            |
+---------------------------+
| ID (PK)                   |
| FirstName                 |
| LastName                  |
| MiddleName                |
| Email (unique)            |
| Password                  |
| PhoneNumber               |
| CreatedAt                 |
| UserType                  |
| DeactivatedAt             |
+---------------------------+

+---------------------------+
|          clients          |
+---------------------------+
| UserID (PK, FK -> users)  |
+---------------------------+

+---------------------------+
|         managers          |
+---------------------------+
| UserID (PK, FK -> users)  |
+---------------------------+

+---------------------------+
|         workers           |
+---------------------------+
| UserID (PK, FK -> users)  |
| ProjectID (unique, FK -> projects) |
+---------------------------+

+---------------------------+
|          projects         |
+---------------------------+
| ID (PK)                   |
| Name                      |
| ClientID (FK -> clients)  |
| ServiceDeposit            |
| MaterialDeposit           |
| CreatedAt                 |
| CompletedAt               |
+---------------------------+

+---------------------------+
|      manager_projects     |
+---------------------------+
| ManagerID (PK, FK -> managers) |
| ProjectID (PK, FK -> projects) |
+---------------------------+

+---------------------------+
|         documents         |
+---------------------------+
| ID (PK)                   |
| Name                      |
| File                      |
| ProjectID (FK -> projects)|
| CreatedAt                 |
+---------------------------+

+---------------------------+
|           acts            |
+---------------------------+
| ID (PK)                   |
| Name                      |
| UnitPrice                 |
| Quantity                  |
| ProjectID (FK -> projects)|
| CreatedAt                 |
+---------------------------+

+---------------------------+
|         services          |
+---------------------------+
| ID (PK)                   |
| Name                      |
| Description               |
| UnitPrice                 |
| CreatedAt                 |
| ArchivedAt                |
+---------------------------+