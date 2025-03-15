CREATE TYPE "states" AS ENUM (
  'proposed',
  'approved',
  'invested',
  'disbursed'
);

CREATE TABLE "borrowers" (
  "id" bigserial PRIMARY KEY,
  "user_name" varchar NOT NULL,
  "email" varchar NOT NULL,
  "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "employees" (
  "id" bigserial PRIMARY KEY,
  "user_name" varchar NOT NULL,
  "email" varchar NOT NULL,
  "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "investors" (
  "id" bigserial PRIMARY KEY,
  "user_name" varchar NOT NULL,
  "email" varchar NOT NULL,
  "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "loans" (
  "id" bigserial PRIMARY KEY,
  "borrower_id" bigint NOT NULL,
  "principal_amount" DECIMAL(20,2),
  "rate" DECIMAL(20,2),
  "roi" DECIMAL(20,2),
  "agreement_date" timestamptz DEFAULT (now()),
  "loan_state" states,
  "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "loan_approvals" (
  "id" bigserial PRIMARY KEY,
  "loan_id" bigint NOT NULL,
  "employee_validator_id" bigint NOT NULL,
  "visit_proof" varchar,
  "approval_date" timestamptz DEFAULT (now()),
  "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "loan_disbursements" (
  "id" bigserial PRIMARY KEY,
  "loan_id" bigint NOT NULL,
  "agreement_letter" varchar,
  "employee_officer_id" bigint NOT NULL,
  "disbursement_date" timestamptz DEFAULT (now()),
  "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "loan_investments" (
  "id" bigserial PRIMARY KEY,
  "loan_id" bigint NOT NULL,
  "investor_id" bigint NOT NULL,
  "agreement_letter" varchar,
  "invested_amount" DECIMAL(20,2),
  "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "loan_states" (
  "id" bigserial PRIMARY KEY,
  "loan_id" bigint NOT NULL,
  "state" states,
  "changed_at" timestamptz DEFAULT (now())
);

CREATE INDEX ON "employees" ("user_name");

CREATE INDEX ON "investors" ("user_name");

CREATE INDEX ON "borrowers" ("user_name");

CREATE INDEX ON "loan_approvals" ("loan_id");

CREATE INDEX ON "loan_disbursements" ("loan_id");

CREATE INDEX ON "loan_investments" ("loan_id");

CREATE INDEX ON "loan_investments" ("investor_id");

CREATE INDEX ON "loan_states" ("loan_id");

ALTER TABLE "loan_approvals" ADD FOREIGN KEY ("loan_id") REFERENCES "loans" ("id");

ALTER TABLE "loans" ADD FOREIGN KEY ("borrower_id") REFERENCES "borrowers" ("id");

ALTER TABLE "loan_approvals" ADD FOREIGN KEY ("employee_validator_id") REFERENCES "employees" ("id");

ALTER TABLE "loan_disbursements" ADD FOREIGN KEY ("loan_id") REFERENCES "loans" ("id");

ALTER TABLE "loan_disbursements" ADD FOREIGN KEY ("employee_officer_id") REFERENCES "employees" ("id");

ALTER TABLE "loan_investments" ADD CONSTRAINT unique_loan_investor UNIQUE ("loan_id", "investor_id");

ALTER TABLE "investors" ADD CONSTRAINT unique_email_investor UNIQUE ("email");

ALTER TABLE "employees" ADD CONSTRAINT unique_email_employee UNIQUE ("email");

ALTER TABLE "loan_investments" ADD FOREIGN KEY ("loan_id") REFERENCES "loans" ("id");

ALTER TABLE "loan_investments" ADD FOREIGN KEY ("investor_id") REFERENCES "investors" ("id");

ALTER TABLE "loan_states" ADD FOREIGN KEY ("loan_id") REFERENCES "loans" ("id");
