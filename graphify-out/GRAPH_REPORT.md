# Graph Report - project_Open  (2026-08-26)

## Corpus Check
- cluster-only mode — file stats not available

## Summary
- 93 nodes · 202 edges · 8 communities (7 shown, 1 thin omitted)
- Extraction: 87% EXTRACTED · 13% INFERRED · 0% AMBIGUOUS · INFERRED: 26 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `8948c2a1`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- Community 0
- Community 1
- Community 2
- Community 3
- Community 4
- Community 5
- Community 6
- Community 7

## God Nodes (most connected - your core abstractions)
1. `run()` - 15 edges
2. `Change` - 11 edges
3. `Resolve()` - 10 edges
4. `Compare()` - 9 edges
5. `comparePath()` - 7 edges
6. `change()` - 7 edges
7. `Run()` - 7 edges
8. `Report()` - 7 edges
9. `resolveCommand()` - 7 edges
10. `Snapshot()` - 6 edges

## Surprising Connections (you probably didn't know these)
- `init()` --calls--> `run()`  [INFERRED]
  main_test.go → main.go
- `TestProbeCommandCapturesResult()` --calls--> `run()`  [INFERRED]
  main_test.go → main.go
- `TestProbeWithoutArgumentsRunsSelfProbe()` --calls--> `run()`  [INFERRED]
  main_test.go → main.go
- `TestReportRedactWritesNewSafeCopy()` --calls--> `run()`  [INFERRED]
  main_test.go → main.go
- `TestResolveAndDiffCommandsProduceJSON()` --calls--> `run()`  [INFERRED]
  main_test.go → main.go

## Import Cycles
- None detected.

## Communities (8 total, 1 thin omitted)

### Community 0 - "Community 0"
Cohesion: 0.27
Nodes (18): Change, Priority, change(), commonEntries(), Compare(), CompareJSON(), comparePath(), CompareProbe() (+10 more)

### Community 1 - "Community 1"
Cohesion: 0.26
Nodes (13): os.FileInfo, candidateInfo(), Result, provenance(), Resolve(), TestResolveCapsReportedCandidatesButKeepsScanning(), TestResolveKeepsExplicitExtension(), TestResolvePOSIXRequiresExecutableFileAndClassifiesProvenance() (+5 more)

### Community 2 - "Community 2"
Cohesion: 0.19
Nodes (9): isWSL(), Snapshot(), split(), Snapshot, TestClassify(), Classify(), Platform, Kind (+1 more)

### Community 3 - "Community 3"
Cohesion: 0.45
Nodes (11): io.Writer, diffCommand(), main(), probeCommand(), reportCommand(), resolveCommand(), run(), snapshotCommand() (+3 more)

### Community 4 - "Community 4"
Cohesion: 0.33
Nodes (10): isSecretName(), markdown(), Report(), sanitizeJSON(), sanitizeText(), TestReportRedactsHomePathsAndRendersJSONMarkdown(), TestReportRedactsStructuredAuthorizationAndGenericKey(), TestTextRedactsNamedAndRecognizedSecrets() (+2 more)

### Community 5 - "Community 5"
Cohesion: 0.31
Nodes (8): context.Context, io.Reader, capture(), environment(), Result, Run(), TestRunCapturesSeparateStreamsExitAndInvalidUTF8(), Capture

### Community 6 - "Community 6"
Cohesion: 0.29
Nodes (8): testing.T, TestSnapshotFixturesDecode(), init(), TestProbeCommandCapturesResult(), TestProbeWithoutArgumentsRunsSelfProbe(), TestReportRedactWritesNewSafeCopy(), TestResolveAndDiffCommandsProduceJSON(), TestSnapshotDoesNotDumpEnvironmentValues()

## Knowledge Gaps
- **2 isolated node(s):** `github.com/TrandomHL/AgentExecTrace`, `Platform`
  These have ≤1 connection - possible missing edges or undocumented components.
- **1 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `run()` connect `Community 3` to `Community 6`?**
  _High betweenness centrality (0.127) - this node is a cross-community bridge._
- **Why does `Report()` connect `Community 4` to `Community 3`?**
  _High betweenness centrality (0.114) - this node is a cross-community bridge._
- **Why does `Run()` connect `Community 5` to `Community 3`?**
  _High betweenness centrality (0.114) - this node is a cross-community bridge._
- **Are the 6 inferred relationships involving `run()` (e.g. with `init()` and `TestProbeCommandCapturesResult()`) actually correct?**
  _`run()` has 6 INFERRED edges - model-reasoned connections that need verification._
- **Are the 5 inferred relationships involving `Resolve()` (e.g. with `TestResolveCapsReportedCandidatesButKeepsScanning()` and `TestResolveKeepsExplicitExtension()`) actually correct?**
  _`Resolve()` has 5 INFERRED edges - model-reasoned connections that need verification._
- **Are the 2 inferred relationships involving `Compare()` (e.g. with `TestCompareReportsSemanticSnapshotChanges()` and `TestGoldenSnapshotPairsDecodeAndCompare()`) actually correct?**
  _`Compare()` has 2 INFERRED edges - model-reasoned connections that need verification._
- **What connects `github.com/TrandomHL/AgentExecTrace`, `Platform` to the rest of the system?**
  _2 weakly-connected nodes found - possible documentation gaps or missing edges._