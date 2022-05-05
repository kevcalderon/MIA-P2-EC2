package main

type MBR struct {
	mbr_tamano         []bytes
	mbr_fecha_creacion []bytes
	mbr_dsk_signatura  []bytes
	dsk_fit            []bytes
	Particiones        [4]Particion
}

type Partion struct {
	part_status []bytes
	part_type   []bytes
	part_fit    []bytes
	part_start  []bytes
	part_size   []bytes
	part_name   []bytes
}

type EBR struct {
	part_status []bytes
	part_fit    []bytes
	part_start  []bytes
	part_size   []bytes
	part_next   []bytes
	part_name   []bytes
}

type Superbloque struct {
	s_filesystem_type   []bytes
	s_inodes_count      []bytes
	s_blocks_count      []bytes
	s_free_blocks_count []bytes
	s_free_inodes_count []bytes
	s_mtime             []bytes
	s_mnt_count         []bytes
	s_magic             []bytes
	s_inode_size        []bytes
	s_block_size        []bytes
	s_firts_ino         []bytes
	s_firts_blo         []bytes
	s_bm_inode_start    []bytes
	s_bm_block_start    []bytes
	s_inode_start       []bytes
	s_block_start       []bytes
}

type inodos struct {
	i_uid   []bytes
	i_gid   []bytes
	i_size  []bytes
	i_atime []bytes
	i_ctime []bytes
	i_mtime []bytes
	i_block []bytes
	i_type  []bytes
	i_perm  []bytes
}

type BlockFile struct {
	b_content []bytes
}
