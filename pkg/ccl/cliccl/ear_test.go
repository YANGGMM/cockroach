// Copyright 2022 The Cockroach Authors.
//
// Licensed as a CockroachDB Enterprise file under the Cockroach Community
// License (the "License"); you may not use this file except in compliance with
// the License. You may obtain a copy of the License at
//
//     https://github.com/cockroachdb/cockroach/blob/master/licenses/CCL.txt

package cliccl

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cockroachdb/cockroach/pkg/ccl/baseccl"
	// The following import is required for the hook that populates
	// NewEncryptedEnvFunc in `pkg/storage`.
	_ "github.com/cockroachdb/cockroach/pkg/ccl/storageccl/engineccl"
	"github.com/cockroachdb/cockroach/pkg/cli"
	"github.com/cockroachdb/cockroach/pkg/settings/cluster"
	"github.com/cockroachdb/cockroach/pkg/storage"
	"github.com/cockroachdb/cockroach/pkg/util/envutil"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/cockroach/pkg/util/randutil"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func TestDecrypt(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	ctx := context.Background()
	dir := t.TempDir()

	// Generate a new encryption key to use.
	keyPath := filepath.Join(dir, "aes.key")
	err := cli.GenEncryptionKeyCmd.RunE(nil, []string{keyPath})
	require.NoError(t, err)

	// Spin up a new encrypted store.
	encSpecStr := fmt.Sprintf("path=%s,key=%s,old-key=plain", dir, keyPath)
	encSpec, err := baseccl.NewStoreEncryptionSpec(encSpecStr)
	require.NoError(t, err)
	encOpts, err := encSpec.ToEncryptionOptions()
	require.NoError(t, err)
	p, err := storage.Open(ctx, storage.Filesystem(dir), cluster.MakeClusterSettings(), storage.EncryptionAtRest(encOpts))
	require.NoError(t, err)

	// Find a manifest file to check.
	files, err := p.List(dir)
	require.NoError(t, err)
	var manifestPath string
	for _, basename := range files {
		if strings.HasPrefix(basename, "MANIFEST-") {
			manifestPath = filepath.Join(dir, basename)
			break
		}
	}
	// Should have found a manifest file.
	require.NotEmpty(t, manifestPath)

	// Close the DB.
	p.Close()

	// Pluck the `pebble manifest dump` command out of the debug command.
	dumpCmd := getTool(cli.DebugPebbleCmd, []string{"pebble", "manifest", "dump"})
	require.NotNil(t, dumpCmd)

	dumpManifest := func(cmd *cobra.Command, path string) string {
		var b bytes.Buffer
		dumpCmd.SetOut(&b)
		dumpCmd.SetErr(&b)
		dumpCmd.Run(cmd, []string{path})
		return b.String()
	}
	out := dumpManifest(dumpCmd, manifestPath)
	// Check for the presence of the comparator line in the manifest dump, as a
	// litmus test for the manifest file being readable. This line should only
	// appear once the file has been decrypted.
	const checkStr = "comparer:     cockroach_comparator"
	require.NotContains(t, out, checkStr)

	// Decrypt the manifest file.
	outPath := filepath.Join(dir, "manifest.plain")
	decryptCmd := getTool(cli.DebugCmd, []string{"debug", "encryption-decrypt"})
	require.NotNil(t, decryptCmd)
	err = decryptCmd.Flags().Set("enterprise-encryption", encSpecStr)
	require.NoError(t, err)
	err = decryptCmd.RunE(decryptCmd, []string{dir, manifestPath, outPath})
	require.NoError(t, err)

	// Check that the decrypted manifest file can now be read.
	out = dumpManifest(dumpCmd, outPath)
	require.Contains(t, out, checkStr)
}

func TestList(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	// Pin the random generator to use a fixed seed. This also requires overriding
	// the PRNG for the duration of the test. This ensures that all the keys
	// generated by the key manager are deterministic.
	reset := envutil.TestSetEnv(t, "COCKROACH_RANDOM_SEED", "1665612120123601000")
	defer reset()
	randBefore := rand.Reader
	randOverride, _ := randutil.NewPseudoRand()
	rand.Reader = randOverride
	defer func() { rand.Reader = randBefore }()

	ctx := context.Background()
	dir := t.TempDir()

	// Generate a new encryption key to use.
	keyPath := filepath.Join(dir, "aes.key")
	err := cli.GenEncryptionKeyCmd.RunE(nil, []string{keyPath})
	require.NoError(t, err)

	// Spin up a new encrypted store.
	encSpecStr := fmt.Sprintf("path=%s,key=%s,old-key=plain", dir, keyPath)
	encSpec, err := baseccl.NewStoreEncryptionSpec(encSpecStr)
	require.NoError(t, err)
	encOpts, err := encSpec.ToEncryptionOptions()
	require.NoError(t, err)
	p, err := storage.Open(ctx, storage.Filesystem(dir), cluster.MakeClusterSettings(), storage.EncryptionAtRest(encOpts))
	require.NoError(t, err)

	// Write a key and flush, to create a table in the store.
	err = p.PutUnversioned([]byte("foo"), nil)
	require.NoError(t, err)
	err = p.Flush()
	require.NoError(t, err)
	p.Close()

	// List the files in the registry.
	cmd := getTool(cli.DebugCmd, []string{"debug", "encryption-registry-list"})
	require.NotNil(t, cmd)
	var b bytes.Buffer
	cmd.SetOut(&b)
	cmd.SetErr(&b)
	err = runList(cmd, []string{dir})
	require.NoError(t, err)

	const want = `000002.log:
  env type: Data, AES128_CTR
  keyID: bbb65a9d114c2a18740f27b6933b74f61018bd5adf545c153b48ffe6473336ef
  nonce: 06 c2 26 f9 68 f0 fc ff b9 e7 82 8f
  counter: 914487965
000004.log:
  env type: Data, AES128_CTR
  keyID: bbb65a9d114c2a18740f27b6933b74f61018bd5adf545c153b48ffe6473336ef
  nonce: 80 18 c0 79 61 c7 cf ef b4 25 4e 78
  counter: 1483615076
000005.sst:
  env type: Data, AES128_CTR
  keyID: bbb65a9d114c2a18740f27b6933b74f61018bd5adf545c153b48ffe6473336ef
  nonce: 71 12 f7 22 9a fb 90 24 4e 58 27 01
  counter: 3082989236
COCKROACHDB_DATA_KEYS_000001_monolith:
  env type: Store, AES128_CTR
  keyID: f594229216d81add7811c4360212eb7629b578ef4eab6e5d05679b3c5de48867
  nonce: 8f 4c ba 1a a3 4f db 3c db 84 cf f5
  counter: 2436226951
CURRENT:
  env type: Data, AES128_CTR
  keyID: bbb65a9d114c2a18740f27b6933b74f61018bd5adf545c153b48ffe6473336ef
  nonce: 18 c2 a6 23 cc 6e 2e 7c 8e bf 84 77
  counter: 3159373900
MANIFEST-000001:
  env type: Data, AES128_CTR
  keyID: bbb65a9d114c2a18740f27b6933b74f61018bd5adf545c153b48ffe6473336ef
  nonce: 2e fd 49 2f 5f c5 53 0a e8 8f 78 cc
  counter: 110434741
OPTIONS-000003:
  env type: Data, AES128_CTR
  keyID: bbb65a9d114c2a18740f27b6933b74f61018bd5adf545c153b48ffe6473336ef
  nonce: d3 97 11 b3 1a ed 22 2b 74 fb 02 0c
  counter: 1229228536
marker.datakeys.000001.COCKROACHDB_DATA_KEYS_000001_monolith:
  env type: Store, AES128_CTR
  keyID: f594229216d81add7811c4360212eb7629b578ef4eab6e5d05679b3c5de48867
  nonce: 55 d7 d4 27 6c 97 9b dd f1 5d 40 c8
  counter: 467030050
`
	require.Equal(t, want, b.String())
}

// getTool traverses the given cobra.Command recursively, searching for a tool
// matching the given command.
func getTool(cmd *cobra.Command, want []string) *cobra.Command {
	// Base cases.
	if cmd.Name() != want[0] {
		return nil
	}
	if len(want) == 1 {
		return cmd
	}
	// Recursive case.
	for _, subCmd := range cmd.Commands() {
		found := getTool(subCmd, want[1:])
		if found != nil {
			return found
		}
	}
	return nil
}
