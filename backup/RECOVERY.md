# Artifact Recovery

Reviewed all 52 currently tracked original JPEGs in `artifacts`, including
`33-old.jpg`, excluding `*-550.jpg`. Inspected 178 distinct image versions
across 80 artifact-related commits from all locally available Git refs.
The repository is not a shallow clone. No remote fetch was performed.

## Results

- 31 images with speaker imagery were saved in this directory.
- 22 were restored from older versions with more speaker images.
- 9 already had suitable current versions and were preserved unchanged.
- 21 had no version containing speaker portraits; these were subsequently copied
  from the current `artifacts` originals as fallbacks, without overwriting any
  existing backup images.
- This directory now contains all 52 originals: events 1 through 51 plus `33-old.jpg`.
- Original files in `artifacts` were not changed. Nothing was committed or pushed.

Recovered images are exact copies of Git blobs, without resizing or re-encoding.
Selection favored the most populated speaker-image slots for the actual event,
then the newest suitable version. Placeholder landscapes do not count as speaker
images. Both populated slots do not necessarily mean every named co-speaker has
an individual portrait.

See [recovery-manifest.json](recovery-manifest.json) for the result of every file,
the number of distinct versions checked, source commit snapshots, Git blob hashes,
image dimensions, and populated speaker slots. Snapshot dates identify a commit
containing the selected image, not necessarily the commit that introduced it.
The manifest records the initial speaker-image recovery; entries marked
`no-speaker-image-found` now have current-original fallback copies in this folder.

## No Speaker Portrait Version Found

| Original | Distinct Versions Checked |
| --- | ---: |
| [1.jpg](../artifacts/1.jpg) | 2 |
| [2.jpg](../artifacts/2.jpg) | 2 |
| [3.jpg](../artifacts/3.jpg) | 3 |
| [4.jpg](../artifacts/4.jpg) | 2 |
| [5.jpg](../artifacts/5.jpg) | 2 |
| [6.jpg](../artifacts/6.jpg) | 2 |
| [7.jpg](../artifacts/7.jpg) | 2 |
| [8.jpg](../artifacts/8.jpg) | 2 |
| [9.jpg](../artifacts/9.jpg) | 2 |
| [10.jpg](../artifacts/10.jpg) | 2 |
| [11.jpg](../artifacts/11.jpg) | 3 |
| [12.jpg](../artifacts/12.jpg) | 2 |
| [13.jpg](../artifacts/13.jpg) | 2 |
| [14.jpg](../artifacts/14.jpg) | 2 |
| [16.jpg](../artifacts/16.jpg) | 2 |
| [17.jpg](../artifacts/17.jpg) | 2 |
| [18.jpg](../artifacts/18.jpg) | 2 |
| [19.jpg](../artifacts/19.jpg) | 2 |
| [33-old.jpg](../artifacts/33-old.jpg) | 1 |
| [33.jpg](../artifacts/33.jpg) | 2 |
| [51.jpg](../artifacts/51.jpg) | 3 |

Event 51 is a pub quiz: its current custom banner has no speaker portraits,
and its older generated versions contain placeholders.

## Partial Recoveries

These are the best versions found, but one speaker slot remains a placeholder.

| Backup | Portrait Present | Still Missing |
| --- | --- | --- |
| [15.jpg](15.jpg) | First slot | Second slot |
| [23.jpg](23.jpg) | Second slot | First slot |
| [28.jpg](28.jpg) | Second slot | First slot |
| [29.jpg](29.jpg) | First slot | Second slot |
| [31.jpg](31.jpg) | Second slot | First slot |
| [32.jpg](32.jpg) | Second slot | First slot |
| [38.jpg](38.jpg) | First slot | Second slot |

## Special Cases

- [42.jpg](42.jpg) has one portrait; the second slot is a community open mic.
  An unrelated older banner titled "My awesome test event" had two portraits
  but was excluded because it does not represent the actual event.
- [45.jpg](45.jpg) uses a cartoon avatar for the first speaker and a photographic
  portrait for the second. No photographic first-slot alternative was found.
- Historical-only deleted derivative sizes such as `*-380.jpg` and `*-600.jpg`,
  and the old `output.jpg` path, were outside the current-file recovery scope.

## Verification

All 31 output hashes match their recorded source Git blobs and source commit
paths. All 52 originals are accounted for in the manifest. The checked-out
commit and tracked working-tree files remained unchanged.

After adding the 21 fallback copies, SHA-256 checks confirmed that they match
their current originals, all 31 existing backup images were unchanged, and no
source images were modified. Every original now has a corresponding backup;
`*-550.jpg` files remain excluded.