- [x] Postgres schema
  - [x] with migrations
  - [ ] what can I use instead of postgres passwords?



- [ ] Client server
  - [x] static upload web page
    - [x] successful upload flash
    - [x] 404 page
    - [x] upload a single file
    - [ ] upload multiple files
    - [ ] upload a directory
    - [x] optional feature: download files if allowed by admin
  - [x] upload endpoint
- [ ] Admin server
  - [x] Link management page
    - [ ] with link management endpoints:
      - [x] create
        - [ ] set maximum file size
        - [ ] DIFFICULT: set maximum total file size
      - [x] edit
      - [x] delete
  - [x] Link view page (allows downloading files)
    - [x] Download file endpoint
  - [ ] optional: multiple admins (a good case to practice migrations!)


- [ ] Tests


File storage:
- [x] $PREFIX is configurable and tells the server where to dump files to
- [x] link:filename must be unique and user cannot reupload the file
- all link:filename files are stored at $PREFIX/$link/$filename
