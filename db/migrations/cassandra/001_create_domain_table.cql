-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS domain (
    GlobalRank int,
    TldRank int,
    Domain text,
    TLD text,
    RefSubNets int,
    RefIPs int,
    IDN_Domain text,
    IDN_TLD text,
    PrevGlobalRank int,
    PrevTldRank int,
    PrevRefSubNets int,
    PrevRefIPs int,
    PRIMARY KEY (Domain)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS domain;
-- +goose StatementEnd
