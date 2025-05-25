'use client';

import { useState } from 'react';
import { Search } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { userApi, ApiError } from '@/lib/api';
import { User } from '@/lib/types';

export function SearchUsers() {
  const [searchType, setSearchType] = useState<'user_id' | 'email'>('user_id');
  const [searchValue, setSearchValue] = useState('');
  const [searchResult, setSearchResult] = useState<User | null>(null);
  const [searching, setSearching] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSearch = async () => {
    if (!searchValue.trim()) return;

    try {
      setSearching(true);
      setError(null);
      setSearchResult(null);

      let user: User;
      if (searchType === 'user_id') {
        user = await userApi.getUserByUserId(searchValue.trim());
      } else {
        user = await userApi.getUserByEmail(searchValue.trim());
      }

      setSearchResult(user);
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message);
      } else {
        setError('Search failed');
      }
    } finally {
      setSearching(false);
    }
  };

  const clearSearch = () => {
    setSearchResult(null);
    setError(null);
    setSearchValue('');
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Search Users</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          <div className="flex gap-4">
            <div className="flex-1">
              <Label htmlFor="search-input">
                Search by {searchType === 'user_id' ? 'User ID' : 'Email'}
              </Label>
              <Input
                id="search-input"
                type={searchType === 'email' ? 'email' : 'text'}
                placeholder={`Enter ${searchType === 'user_id' ? 'user ID' : 'email address'}`}
                value={searchValue}
                onChange={(e) => setSearchValue(e.target.value)}
                onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
              />
            </div>
            <div className="flex flex-col gap-2">
              <Label>Search Type</Label>
              <div className="flex gap-2">
                <Button
                  variant={searchType === 'user_id' ? 'default' : 'outline'}
                  size="sm"
                  onClick={() => setSearchType('user_id')}
                >
                  User ID
                </Button>
                <Button
                  variant={searchType === 'email' ? 'default' : 'outline'}
                  size="sm"
                  onClick={() => setSearchType('email')}
                >
                  Email
                </Button>
              </div>
            </div>
          </div>

          <div className="flex gap-2">
            <Button
              onClick={handleSearch}
              disabled={searching || !searchValue.trim()}
            >
              <Search className="mr-2 h-4 w-4" />
              {searching ? 'Searching...' : 'Search'}
            </Button>
            {(searchResult || error) && (
              <Button variant="outline" onClick={clearSearch}>
                Clear
              </Button>
            )}
          </div>

          {error && (
            <div className="p-3 rounded-md bg-destructive/10 border border-destructive/20">
              <p className="text-sm text-destructive">{error}</p>
            </div>
          )}

          {searchResult && (
            <div className="p-4 rounded-md bg-muted/50 border">
              <h4 className="font-medium mb-2">Search Result:</h4>
              <div className="grid grid-cols-2 gap-2 text-sm">
                <div><strong>User ID:</strong> {searchResult.user_id}</div>
                <div><strong>Email:</strong> {searchResult.email}</div>
                <div><strong>MongoDB ID:</strong> {searchResult.id}</div>
                <div><strong>Created:</strong> {new Date(searchResult.created_at).toLocaleString()}</div>
              </div>
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
}