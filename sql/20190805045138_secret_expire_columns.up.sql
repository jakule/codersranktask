ALTER TABLE secrets
    ADD expireAfterViews INT;

ALTER TABLE secrets
    ADD expireAfterTime TIMESTAMP;