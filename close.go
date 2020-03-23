package redis

// Close closes connection pool.
func (r *Client) Close() error {
	err := r.pool.Close()
	return err
}