module.exports = {
  branches: ['main'],
  plugins: [
    [
      '@semantic-release/commit-analyzer',
      { preset: 'conventionalcommits' }
    ],
    [
      '@semantic-release/release-notes-generator',
      {
        preset: 'conventionalcommits',
        releaseNotesGeneratorOpts: {
          transform: (commit) => {
            if (!commit.committerDate || isNaN(new Date(commit.committerDate).getTime())) {
              return null;
            }
            return commit;
          }
        }
      }
    ],
    '@semantic-release/changelog',
    '@semantic-release/github',
    [
      '@semantic-release/git',
      {
        assets: ['CHANGELOG.md'],
        message: 'chore(release): ${nextRelease.version} [skip ci]\n\n${nextRelease.notes}'
      }
    ]
  ]
};
