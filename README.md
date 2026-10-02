# Shipping Rate Calculator

Command-line program that calculates the shipping cost of a parcel from its
weight and destination zone.

## Running

Requires Go 1.22 or newer. There are no third-party dependencies.

```sh
go run ./cmd/shiprate -weight 3 -length 30 -width 20 -height 10 -zone eu
```

```
weight tier        10.00
zone surcharge      8.00
total              18.00 EUR
```

| Flag      | Meaning                             |
|-----------|-------------------------------------|
| `-weight` | weight in kg                        |
| `-length` | length in cm                        |
| `-width`  | width in cm                         |
| `-height` | height in cm                        |
| `-zone`   | `domestic`, `eu` or `international` |

All flags are required. Invalid input prints the reason to stderr and exits
with status 1.

## Tests

```sh
go test ./...
```

Each test file sits next to the code it tests. The ones in `internal/shipping`
check the prices: every weight tier with every zone, the 1 kg and 5 kg limits,
and invalid parcels. The ones in `cmd/shiprate` run the whole program and check
what it prints.

## Design

`cmd/shiprate` is the program: it reads the flags and prints the result. The
prices are worked out in `internal/shipping`.

A parcel goes through a list of rules. Each rule returns one charge, and the
total is the sum of the charges. Today there are two rules, the weight tier and
the zone surcharge. Their prices are in `rates.go`.

To change a price or add a tier, edit `rates.go`. A new zone is added to the
list of zones in `parcel.go` and gets its surcharge in `rates.go`. A new kind
of charge is a new type with an `Apply` method, added to the list in
`rates.go`.

Rules do not see each other's charges, so a percentage discount on the subtotal
would need a change to the `Rule` interface.

## Assumptions

- A weight or dimension that is zero or negative is rejected with an error.
- Length, width and height are required because a parcel has them, but they do
  not change the price.
- The zone can be typed in any letter case (`EU`, `eu`).
- Amounts are kept in integer cents.
