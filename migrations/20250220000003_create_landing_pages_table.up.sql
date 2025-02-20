CREATE TABLE landing_pages (
    id          varchar(255) PRIMARY KEY,
    slug        varchar(55) NOT NULL REFERENCES tracking_links(slug),
    title       varchar(255) NOT NULL,
    preview_url text,
    target_url  text NOT NULL,
    created_at  timestamp without time zone DEFAULT NOW(),
    updated_at  timestamp without time zone DEFAULT NOW()
);

CREATE INDEX idx_landing_pages_slug ON landing_pages(slug);