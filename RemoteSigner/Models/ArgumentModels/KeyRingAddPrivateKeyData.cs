﻿using System;
namespace RemoteSigner.Models.ArgumentModels {
    public struct KeyRingAddPrivateKeyData {
        public String EncryptedPrivateKey { get; set; }
        public bool SaveToDisk { get; set; }
    }
}
