BEGIN TRANSACTION;

DROP TABLE IF EXISTS ANNOUNCEMENTS; 
CREATE TABLE ANNOUNCEMENTS (
    ID INTEGER PRIMARY KEY,
    COMBINED_HASH BLOB(32) NOT NULL
);

COMMIT;