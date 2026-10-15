# sparse-equation-system-solve

A toolkit for the decomposition of network-type sparse systems of linear algebraic equations (SLAE).

## Decomposition Library Overview (`internal/solver`)

The mathematical core of the library implements a decomposition method for rectangular underdetermined SLAEs of the form Ax = b on graphs.

## Benchmark Task (`cmd/cloud-net-cost-optimization-task`)

A CLI application that solves a **quadratic programming** problem—distributing data flows in a cloud server cluster to achieve minimal operational costs. The execution compares a standard solution approach with a graph-based decomposition alternative.

## Web Application (`web/app`)

A backend-less client web application (SPA) for interactive decomposition and solving of algebraic equation systems directly in the browser.

- **Top Section** - A graph configuration form. It allows users to dynamically change the number of vertices and edges, as well as select the tree root for the decomposition algorithm. It includes automatic input validation for vector b and safeguards against graph loops before computations start.
- **Bottom Section**:
    1. *Graphical Block* - Visualization of the graph and the discovered spanning tree.
    2. *Mathematical Block* - A comprehensive report on the system solution (the particular solution vector, the fundamental system of solutions (FSS), and the final general parametric solution system).

## Tech Stack

- Go 1.22+ SDK, `gonum.org/v1/gonum`, Go WASM module,
- React 19, Vite 8, TypeScript, Apache ECharts (graph visualization), KaTeX (LaTeX syntax rendering support),
- GNU Make toolchain

## Project Commands

All core build, test, and deployment operations are managed via the `Makefile` (supported on Linux, and Windows via Git Bash or WSL) in the root directory:

- **Run Go test suite**
```bash
make test
```
- **Run benchmark utility**
```bash
go run ./cmd/cloud-net-cost-optimization-task/
```
- **Compile WASM module and build production frontend:**
```bash
make build
```
- **Launch local server to debug production build:**
```bash
cd web/app && npm run preview
```
