CREATE TABLE redirect_rules (
    id                  serial PRIMARY KEY,
    redirect_type       varchar(50) NOT NULL,
    redirect_slug       varchar(55),
    redirect_url        text,
    redirect_smart_slug jsonb,
    created_at         timestamp without time zone DEFAULT NOW(),
    updated_at         timestamp without time zone DEFAULT NOW()
);

-- Add indexes for commonly queried fields
CREATE INDEX idx_redirect_rules_redirect_slug ON redirect_rules(redirect_slug);
CREATE INDEX idx_redirect_rules_redirect_smart_slug ON redirect_rules(redirect_smart_slug);