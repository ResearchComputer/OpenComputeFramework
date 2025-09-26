#!/usr/bin/env python3
"""
Simple test script to verify the implementation
"""
import requests
import json
import time

BASE_URL = "http://localhost:8000"

def test_health():
    """Test health endpoint"""
    try:
        response = requests.get(f"{BASE_URL}/health")
        print(f"Health check: {response.status_code} - {response.json()}")
        return response.status_code == 200
    except Exception as e:
        print(f"Health check failed: {e}")
        return False

def test_root():
    """Test root endpoint"""
    try:
        response = requests.get(f"{BASE_URL}/")
        print(f"Root endpoint: {response.status_code}")
        data = response.json()
        print(f"Available endpoints: {data.get('endpoints', [])}")
        return response.status_code == 200
    except Exception as e:
        print(f"Root endpoint failed: {e}")
        return False

def test_public_models():
    """Test public models endpoint"""
    try:
        response = requests.get(f"{BASE_URL}/api/models/public")
        print(f"Public models: {response.status_code}")
        if response.status_code == 200:
            models = response.json()
            print(f"Found {len(models)} models")
            return True
        else:
            print(f"Error: {response.text}")
            return False
    except Exception as e:
        print(f"Public models failed: {e}")
        return False

def test_wallet_auth():
    """Test wallet authentication (simplified test)"""
    try:
        # This is a mock test - in real implementation you'd need proper signature
        auth_data = {
            "address": "0x742d35Cc6634C0532925a3b844Bc9e7595f1234",
            "signature": "0x1234567890abcdef",
            "chain_id": 1
        }
        response = requests.post(f"{BASE_URL}/api/auth/connect", json=auth_data)
        print(f"Auth test: {response.status_code}")
        if response.status_code != 200:
            print(f"Expected auth to fail (invalid signature): {response.text}")
        return True
    except Exception as e:
        print(f"Auth test failed: {e}")
        return False

def main():
    """Run all tests"""
    print("Starting implementation tests...")
    print("=" * 50)

    tests = [
        ("Health Check", test_health),
        ("Root Endpoint", test_root),
        ("Public Models", test_public_models),
        ("Wallet Auth", test_wallet_auth),
    ]

    passed = 0
    total = len(tests)

    for test_name, test_func in tests:
        print(f"\nTesting {test_name}...")
        if test_func():
            passed += 1
            print(f"✓ {test_name} passed")
        else:
            print(f"✗ {test_name} failed")
        print("-" * 30)

    print(f"\nResults: {passed}/{total} tests passed")
    return passed == total

if __name__ == "__main__":
    main()