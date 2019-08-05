ALTER TABLE secrets
    ADD createdTime timestamp DEFAULT now() NOT NULL;