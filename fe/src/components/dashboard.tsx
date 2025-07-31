import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import {
    Typography,
    List,
    ListItemButton,
    ListItemText,
    CircularProgress,
    Card,
    CardContent,
    Box,
    AppBar,
    Toolbar,
    Button,
    Dialog,
    DialogTitle,
    DialogContent,
    TextField,
    DialogActions,
    IconButton,
    Tooltip,
} from "@mui/material";
import LockIcon from "@mui/icons-material/Lock";
import LockOpenIcon from "@mui/icons-material/LockOpen";
import DeleteIcon from "@mui/icons-material/Delete";
import { getBlocks, createBlock, toggleLock, deleteBlock } from "../api/api";

export default function Dashboard() {
    const [blocks, setBlocks] = useState<any[]>([]);
    const [loading, setLoading] = useState(true);
    const [open, setOpen] = useState(false);
    const [month, setMonth] = useState("");
    const [members, setMembers] = useState("");
    const username = localStorage.getItem("username") || "User";
    const navigate = useNavigate();

    const fetchBlocks = () => {
        getBlocks()
            .then((res) => {
                setBlocks(res.data && Array.isArray(res.data) ? res.data : []);
            })
            .catch((err) => {
                console.error("Failed to load blocks", err);
                setBlocks([]);
                if (err.status === 401) {
                    localStorage.removeItem("token");
                    localStorage.removeItem("username");
                    window.location.href = "/login";
                }
            })
            .finally(() => setLoading(false));
    };

    useEffect(() => {
        fetchBlocks();
    }, []);

    const handleCreateBlock = () => {
        createBlock(month, members.split(","))
            .then(() => {
                setOpen(false);
                setMonth("");
                setMembers("");
                fetchBlocks();
            })
            .catch((err) => console.error("Failed to create block", err));
    };

    const handleToggleLock = (month: string, locked: boolean) => {
        toggleLock(month, locked)
            .then(fetchBlocks)
            .catch(err => console.error("Failed to toggle lock", err));
    };

    const handleDeleteBlock = (month: string) => {
        if (window.confirm(`Bạn có chắc chắn muốn xóa tháng ${month} không?`)) {
            deleteBlock(month)
                .then(fetchBlocks)
                .catch(err => console.error("Failed to delete block", err));
        }
    };

    if (loading)
        return (
            <Box display="flex" justifyContent="center" alignItems="center" height="100vh">
                <CircularProgress />
            </Box>
        );

    return (
        <>
            <AppBar position="static" sx={{ backgroundColor: "#1976d2" }}>
                <Toolbar sx={{ display: "flex", justifyContent: "space-between" }}>
                    <Typography variant="h6">Home</Typography>
                    <Box display="flex" alignItems="center" gap={2}>
                        <Typography variant="subtitle1" component="span">
                            Chào, {username}
                        </Typography>
                        <Button
                            color="inherit"
                            onClick={() => {
                                localStorage.removeItem("token");
                                window.location.href = "/login";
                            }}
                        >
                            Đăng xuất
                        </Button>
                    </Box>
                </Toolbar>
            </AppBar>

            <Box
                sx={{
                    minHeight: "100vh",
                    backgroundColor: "#f5f6fa",
                    padding: "20px",
                    display: "flex",
                    justifyContent: "center",
                }}
            >
                <Card sx={{ width: "100%", maxWidth: "900px", boxShadow: 3, borderRadius: 3 }}>
                    <CardContent>
                        <Typography variant="h5" gutterBottom align="center">
                            Danh sách các tháng hiện tại
                        </Typography>
                        <List>
                            {blocks.length > 0 ? (
                                blocks.map((block, i) => (
                                    <ListItemButton
                                        key={i}
                                        sx={{
                                            borderBottom: "1px solid #f0f0f0",
                                            "&:hover": { backgroundColor: "#f9f9f9" },
                                            display: "flex",
                                            justifyContent: "space-between",
                                        }}
                                        onClick={() => navigate(`/blocks/${block.month}`)}
                                    >
                                        <ListItemText
                                            primary={`${i + 1}. Tháng: ${block.month}`}
                                            secondary={block.locked ? "Đã khóa" : "Chưa khóa"}
                                        />
                                        <Box display="flex" alignItems="center" gap={1}>
                                            <Tooltip title={block.locked ? "Mở khóa tháng này" : "Khóa tháng này"}>
                                                <IconButton
                                                    edge="end"
                                                    onClick={(e) => {
                                                        e.stopPropagation();
                                                        handleToggleLock(block.month, block.locked);
                                                    }}
                                                >
                                                    {block.locked ? (
                                                        <LockIcon color="error" />
                                                    ) : (
                                                        <LockOpenIcon color="success" />
                                                    )}
                                                </IconButton>
                                            </Tooltip>
                                            <Tooltip title="Xóa tháng này">
                                                <IconButton
                                                    edge="end"
                                                    onClick={(e) => {
                                                        e.stopPropagation();
                                                        handleDeleteBlock(block.id);
                                                    }}
                                                >
                                                    <DeleteIcon color="error" />
                                                </IconButton>
                                            </Tooltip>
                                        </Box>
                                    </ListItemButton>
                                ))
                            ) : (
                                <Typography
                                    variant="body2"
                                    color="textSecondary"
                                    align="center"
                                    sx={{ mt: 2 }}
                                >
                                    Không có tháng nào được tạo
                                </Typography>
                            )}
                        </List>
                        <Box mt={2} textAlign="center">
                            <Button variant="contained" color="secondary" onClick={() => setOpen(true)}>
                                Tạo tháng mới
                            </Button>
                        </Box>
                    </CardContent>
                </Card>
            </Box>

            <Dialog open={open} onClose={() => setOpen(false)}>
                <DialogTitle>Tạo tháng mới</DialogTitle>
                <DialogContent>
                    <TextField
                        label="Tháng (YYYY-MM)"
                        value={month}
                        onChange={(e) => setMonth(e.target.value)}
                        fullWidth
                        sx={{ mt: 2 }}
                    />
                    <TextField
                        label="Thành viên (cách nhau bởi dấu phẩy)"
                        value={members}
                        onChange={(e) => setMembers(e.target.value)}
                        fullWidth
                        sx={{ mt: 2 }}
                    />
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setOpen(false)}>Hủy</Button>
                    <Button onClick={handleCreateBlock} variant="contained">
                        Tạo
                    </Button>
                </DialogActions>
            </Dialog>
        </>
    );
}
