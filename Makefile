abigen:
	abigen --sol testdata/counter.sol --pkg testdata --type Counter --out testdata/counter.go

.PHONY: abigen
