# SpannerPlayground
My playground for Google Cloud Spanner

## How to play it
Source code is `pagination/main.go`

Run `tilt up`, and click "pagination" tab, then you can see the pagination effect:
```
Running cmd: go run pagination/main.go
SQL for first page: SELECT PublicApplicationID, Name, Time FROM Application ORDER BY Time DESC LIMIT 5
28801695-ec1f-4d27-bb99-286bc430af81 Zhi Qu 19 19
e882f9d0-d740-47f8-ad4e-cfc6c8b7c30a Zhi Qu 18 18
2d3959e5-b56c-43ab-bdde-bc2a493313e4 Zhi Qu 17 17
4980dc05-0149-459f-a8a0-74e3f4c2ad91 Zhi Qu 16 16
3928966e-7217-489f-9bce-98e396f70d38 Zhi Qu 15 15
Query is done. Next page token: dGltZV8xNQ==
SQL for first page: SELECT PublicApplicationID, Name, Time FROM Application WHERE Time < 15 ORDER BY Time DESC LIMIT 5
b6d716b4-c449-411e-93b1-e56eb44bceb7 Zhi Qu 14 14
f4c91327-734c-4a85-9340-f4478dcfc6b3 Zhi Qu 13 13
69705377-7855-4e11-9720-958a838d90ac Zhi Qu 12 12
368c757f-6e2d-40b8-9596-ac8c5e97dfba Zhi Qu 11 11
b1369f45-c27a-4b99-8d99-14860772df36 Zhi Qu 10 10
Query is done. Next page token: dGltZV8xMA==
SQL for first page: SELECT PublicApplicationID, Name, Time FROM Application WHERE Time < 10 ORDER BY Time DESC LIMIT 5
63e9900f-19bc-4446-9b97-c625172e4ebe Zhi Qu 9 9
ca67e53b-97d4-4320-b3c7-7820b8884abd Zhi Qu 8 8
ff35f5a5-2354-4899-a832-f29a7b6a2c86 Zhi Qu 7 7
338a50c0-4bfa-4194-ba39-79e1a8a48d79 Zhi Qu 6 6
5ce60c98-7ce7-43e0-88cc-4f6862bcd3e6 Zhi Qu 5 5
Query is done. Next page token: dGltZV81
SQL for first page: SELECT PublicApplicationID, Name, Time FROM Application WHERE Time < 5 ORDER BY Time DESC LIMIT 5
930dbf8b-490e-4c00-8543-029a2ac60a28 Zhi Qu 4 4
4f9b5628-dd9a-466f-b609-691a8f86c574 Zhi Qu 3 3
29402822-0934-42da-9232-cb28fa2d1dbd Zhi Qu 2 2
07088a5d-0e01-4991-b36b-0926c6efeaae Zhi Qu 1 1
c405940f-d3cb-4675-9a57-9adb5d5f60a3 Zhi Qu 0 0
Query is done. Next page token: dGltZV8w
SQL for first page: SELECT PublicApplicationID, Name, Time FROM Application WHERE Time < 0 ORDER BY Time DESC LIMIT 5
Pagination is done.
```