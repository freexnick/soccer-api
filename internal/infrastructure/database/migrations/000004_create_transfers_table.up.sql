CREATE TABLE IF NOT EXISTS transfer_listings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    player_id UUID NOT NULL UNIQUE,
    selling_team_id UUID NOT NULL,
    asking_price INT NOT NULL CHECK (asking_price > 0),
    listed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_transfer_listings_player
        FOREIGN KEY(player_id)
        REFERENCES players(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_transfer_listings_selling_team
        FOREIGN KEY(selling_team_id)
        REFERENCES teams(id)
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_transfer_listings_player_id ON transfer_listings(player_id);
CREATE INDEX IF NOT EXISTS idx_transfer_listings_selling_team_id ON transfer_listings(selling_team_id);