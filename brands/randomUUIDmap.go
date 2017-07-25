package brands

/*
   This map is a static map of TME identifiers to the randomly generated UUIDs that were created by editorial
   as part of the Bertha spreadsheet.  Since they don't follow the normal manner of generating a UUID in UP,
   we need to store them as a map TME -> UUID.

   This shouldn't require changing as new brands added to the spreadsheet have their UUIDs generated correctly
   from their TME identifier.
*/
func berthaUUIDmap() map[string]string {
	return map[string]string{
		"QnJhbmRzXzE4-QnJhbmRz":                                     "6773e864-78ab-4051-abc2-f4e9ab423ebb",
		"ZDIyZjNlOWQtZDNlMy00Yzk4LTkwOTctNzFiYjc5ZjNjYTkw-QnJhbmRz": "164d0c3b-8a5a-4163-9519-96b57ed159bf",
		"YmE1NjE5NjMtNWUzOC00MGVlLWJjNjEtYjY5NmMzYTA1MDBk-QnJhbmRz": "806d05b8-3d29-4e81-8668-e9cdae0ab086",
		"QnJhbmRzXzg2-QnJhbmRz":                                     "fd4459b2-cc4e-4ec8-9853-c5238eb860fb",
		"QnJhbmRzXzIwOA==-QnJhbmRz":                                 "5e5dba2b-8031-4a65-b193-84a58882c62e",
		"ZDkyYTVhMzYtYjAyOS00OWI1LWI5ZTgtM2QyYTIzYjk4Y2Jj-QnJhbmRz": "89d15f70-640d-11e4-9803-0800200c9a66",
		"N2I3NDY4ODktYTRjMC00NWNkLTgwZDEtNjgyMTgzM2Q2MmZl-QnJhbmRz": "0cfee514-db0b-4474-b7a6-bf4db16a1380",
		"ZDIzYmQzNjgtYzAzMS00MDg3LTg3MTEtMDUyZTZkZTkzMTc4-QnJhbmRz": "3a37a89e-14ce-4ac8-af12-961a9630dce3",
		"NDdkYzA2MzgtNzczNS00OGMxLWFhNjgtMjI0YWU0MzkyNmJh-QnJhbmRz": "cfc05027-a3fa-472a-a14b-fa2506df665f",
		"YmRhODE3ODMtMTYzZi00NDdjLTg5MmUtY2E2ZGQ5ZTA2NDJh-QnJhbmRz": "13006c72-7d1b-47a0-96fe-d1ad1f12de9f",
		"MGY2ZTQ3MTYtYjJiNS00ODVhLTlkYTktNzZlNzc3YTcxOWYy-QnJhbmRz": "b8513403-7892-4901-bb97-1765fc0ba190",
		"MGM0NTAxM2YtOWIwNS00ODhkLWI0NTgtMTgzMjJjNzUyNTVj-QnJhbmRz": "2f5d019b-9aa4-43bb-b204-c7437bf0f031",
		"NTQ2MGRhM2UtMGNlOC00NDBjLTgyNWEtY2VmMWZkMjk5NDdk-QnJhbmRz": "cd5e45bd-707a-47e9-a478-b18c02ef685f",
		"Yjg0NDNlODEtNjc1YS00ZWU1LWE2MzUtOTA5MDIyNjhhMzQw-QnJhbmRz": "8348b79b-5665-409c-8696-6904f9a26fc6",
		"MmYyZGMzZmQtNmMzZi00OTQ1LThlMmQtZDI0OGU2ZDk2N2M0-QnJhbmRz": "b4fac748-a2b1-4b7d-8e1f-03ba743ff717",
		"MjI0Zjk2YzctZmFiYS00YWY2LWEzZjktZjgxOWU0ODE1YWUz-QnJhbmRz": "72349b8e-1cac-45b4-89df-a1afb55b782e",
		"MWM5NjVlYTYtZGUyMy00NjUwLWEwNWMtYTIwMDJjMTcwMTYw-QnJhbmRz": "e363dfb8-f6d9-4f2c-beba-5162b334272b",
		"YzhlNzZkYTctMDJiNy00NTViLTk3NmYtNmJjYTE5NDEyM2Yw-QnJhbmRz": "2d3e16e0-61cb-4322-8aff-3b01c59f4daa",
		"Y2ZjNGRkNmItMDEzYi00ZTgyLWI4YzYtMjMwODA2OWM2NmU5-QnJhbmRz": "e180d0e0-9d13-4212-9696-f86ef197d2bd",
		"NGY0ZGU0MTUtZTczZC00NzNjLWI2NTgtMjAwOTExMDMwZjI0-QnJhbmRz": "462507e2-e20a-431d-9648-a9131770b3aa",
		"NWE4ZjU5ZWYtNThjNS00ZjhjLWE5MzEtZTQ3ZjczNjEzY2Jk-QnJhbmRz": "556f7ab4-474c-40e9-b7a3-c3af60d34156",
		"ODEwNmQ2NGItOTdlNS00MzBiLThjYmMtMjhiMDhmYjU3YjE2-QnJhbmRz": "a98e560a-e5ef-49c4-9af8-2d7419a77c31",
		"N2Y2ZjFlMDgtNjY3OS00NTU5LThjNWMtNWNiYzcwNTQyNmFj-QnJhbmRz": "0b8b1a83-5897-408a-9232-821c62afa231",
		"MDVlM2Q2OGMtOWRkZi00NzI5LWJlNjctNzhlYmE1MGE3YzJk-QnJhbmRz": "2a2c87f8-eb14-48c2-900f-2cffc2d7c168",
		"QnJhbmRzXzEwOQ==-QnJhbmRz":                                 "95d802a1-106a-4393-a143-e39a364b31e6",
		"OTk5ZDY2OTUtZWNlYS00MjA4LWE0ZmUtYmQzNDY4NTUxMTYx-QnJhbmRz": "b7564b93-9809-4ecd-9b82-bd2efaeb282f",
		"OTI1YmZiNzgtN2IzMS00NjVjLThiODgtY2EzMmQwOGZkMDFh-QnJhbmRz": "51d51df0-a824-4ec7-aa33-faa013c2d3f9",
		"MjQyNGVmODUtYTk3NC00ZWU5LWFhMWUtMjhkYjI3M2ZkZTFm-QnJhbmRz": "d56ac153-41e4-4f3d-aec8-14184d7e1048",
		"NGM5MDA1ZWQtZTI3Zi00NWQwLWI3NWItMjMyMGFlNzY4NTFj-QnJhbmRz": "c5f12e29-f178-4b37-8a81-e9029d0c2609",
		"YmQ4Y2VlMjQtOGY0OS00NmZjLWE1NWItYjgxNzJiMjhlMTRh-QnJhbmRz": "51d51df0-a814-4ec7-aa33-faa013c2d3f9",
		"NDY0OWQ2YjgtMTI0NC00MTk4LTg3OGEtNjFjODFmYzZmZWU1-QnJhbmRz": "27142330-e0bb-49c7-91e7-c51648f0ce68",
		"N2NkMjJiYzQtOGI3MC00NTM4LTgzYmYtMTQ3YmJkZGZkODJj-QnJhbmRz": "5c21e52f-48f9-3a77-ad36-318163f284be",
		"MjFjOTI0Y2YtNGFlOS00OTMzLWJhMjEtNjBjNjE2YTRhMmJi-QnJhbmRz": "7c938825-8c72-3458-b484-54c4d95bf7bd",
		"MWM0NzM0NjEtNjQ0Yi00OGEzLTgxZTgtZWY0NzA5ZDM3ZjU2-QnJhbmRz": "6b82551f-f532-3460-aded-07e1dde00103",
		"ZWQzZDQwOWMtNTAzNy00ZjlhLWExZmMtMzRhZjU4ZTVjOWRk-QnJhbmRz": "b11c39db-ce7f-3e2a-ba30-4781ef63491e",
		"YzQ5NjY4NjEtNjFkMS00M2Q2LTlmNmQtYzcxMzBkMzc1MTg4-QnJhbmRz": "ebd0c569-dd8d-3787-b77e-4db976bcfab1",
		"QnJhbmRzXzIwNA==-QnJhbmRz":                                 "ea1da823-3ab2-3b00-805e-2d3c75994e73",
		"ZjJjOTIxN2EtYTg5NS00YTQ0LTkzZmYtZmU2NGM5MTkzNTIx-QnJhbmRz": "a54fda40-7fe7-339a-9b83-2d7b964ff3a4",
		"MTg2ZmNhMGEtZmE4MC00Y2VlLWI0YjItNGQ5Njc1NzQ1YjJk-QnJhbmRz": "01438770-10b4-343d-97cc-8bd6e27a4fd1",
		"MGY3OTBkZGEtYjIyZC00MDkxLWI1Y2EtZWJlNTVkYzk1YjZh-QnJhbmRz": "c24d6335-076a-366a-98e2-500bb26401d6",
		"OGQ0NzZkYTEtZTRjZS00MTNlLTk1MDYtNzFmOWI1YTIxNGNj-QnJhbmRz": "e258a9b2-1049-3171-b55c-16cab615c47c",
		"NTlhNzEyMzMtZjBjZi00Y2U1LTg0ODUtZWVjNmEyYmU1NzQ2-QnJhbmRz": "5c7592a8-1f0c-11e4-b0cb-b2227cce2b54",
		"ZWQzYjZlYzUtNjQ2Ni00N2VmLWIxZDgtMTY5NTJmZDUyMmM3-QnJhbmRz": "ed3b6ec5-6466-47ef-b1d8-16952fd522c7",
	}
}
