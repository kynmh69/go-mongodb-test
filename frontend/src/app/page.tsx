'use client';

import { useState, useEffect } from 'react';
import { Plus } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { UserTable } from '@/components/UserTable';
import { CreateUserDialog } from '@/components/CreateUserDialog';
import { UpdateUserDialog } from '@/components/UpdateUserDialog';
import { DeleteUserDialog } from '@/components/DeleteUserDialog';
import { SearchUsers } from '@/components/SearchUsers';
import { userApi, ApiError } from '@/lib/api';
import { User } from '@/lib/types';

export default function HomePage() {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  
  // Dialog states
  const [createDialogOpen, setCreateDialogOpen] = useState(false);
  const [updateDialogOpen, setUpdateDialogOpen] = useState(false);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [selectedUser, setSelectedUser] = useState<User | null>(null);

  const fetchUsers = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await userApi.getUsers();
      setUsers(response.users || []);
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message);
      } else {
        setError('Failed to fetch users');
      }
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchUsers();
  }, []);

  const handleCreateUser = async () => {
    await fetchUsers();
    setCreateDialogOpen(false);
  };

  const handleUpdateUser = async () => {
    await fetchUsers();
    setUpdateDialogOpen(false);
    setSelectedUser(null);
  };

  const handleDeleteUser = async () => {
    await fetchUsers();
    setDeleteDialogOpen(false);
    setSelectedUser(null);
  };

  const openUpdateDialog = (user: User) => {
    setSelectedUser(user);
    setUpdateDialogOpen(true);
  };

  const openDeleteDialog = (user: User) => {
    setSelectedUser(user);
    setDeleteDialogOpen(true);
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="text-lg">Loading users...</div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-3xl font-bold tracking-tight">Users</h2>
          <p className="text-muted-foreground">
            Manage users in your system
          </p>
        </div>
        <Button onClick={() => setCreateDialogOpen(true)}>
          <Plus className="mr-2 h-4 w-4" />
          Add User
        </Button>
      </div>

      <SearchUsers />

      {error && (
        <Card className="border-destructive">
          <CardHeader>
            <CardTitle className="text-destructive">Error</CardTitle>
          </CardHeader>
          <CardContent>
            <p>{error}</p>
            <Button
              variant="outline"
              onClick={fetchUsers}
              className="mt-4"
            >
              Retry
            </Button>
          </CardContent>
        </Card>
      )}

      <UserTable
        users={users}
        onEdit={openUpdateDialog}
        onDelete={openDeleteDialog}
      />

      <CreateUserDialog
        open={createDialogOpen}
        onOpenChange={setCreateDialogOpen}
        onSuccess={handleCreateUser}
      />

      {selectedUser && (
        <>
          <UpdateUserDialog
            open={updateDialogOpen}
            onOpenChange={setUpdateDialogOpen}
            user={selectedUser}
            onSuccess={handleUpdateUser}
          />
          <DeleteUserDialog
            open={deleteDialogOpen}
            onOpenChange={setDeleteDialogOpen}
            user={selectedUser}
            onSuccess={handleDeleteUser}
          />
        </>
      )}
    </div>
  );
}