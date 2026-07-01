# Import Cycle Fix

## Problem

The project had a circular dependency between the `normalizer` and `providers/mpesa` packages:

```text
providers/mpesa
        ↓
normalizer
        ↓
providers/mpesa
```

This happened because:

* `providers/mpesa` imported `normalizer` to use the `NormalizedTransaction` type.
* `normalizer` imported `providers/mpesa` to register the M-Pesa normalizer.

Go does not allow cyclic imports, so the project failed to compile.

## Solution

The shared `NormalizedTransaction` model was moved into its own package:

```text
internal/transactions
```

Now:

* `providers/mpesa` imports `transactions`.
* `normalizer` imports both `transactions` and `providers/mpesa`.
* `transactions` has no dependencies.

The dependency graph is now one-way:

```text
transactions
      ↑
      │
providers/mpesa
      ↑
      │
normalizer
```

This removes the circular dependency and creates a cleaner architecture where shared models are independent of the packages that use them.
