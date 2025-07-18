import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import axios from "axios";
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

export default function Dashboard() {
  const [blocks, setBlocks] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [month, setMonth] = useState("");
  const [members, setMembers] = useState("");
  const username = localStorage.getItem("username") || "User";
  const token = localStorage.getItem("token");
  const navigate = useNavigate();

  const fetchBlocks = () => {
    axios
      .get("http://localhost:3000/blocks", {
        headers: { Authorization: `Bearer ${token}` },
      })
      .then((res) => {
        setBlocks(res.data && Array.isArray(res.data) ? res.data : []);
      })
      .catch((err) => {
        console.error("Failed to load blocks", err);
        setBlocks([]);
      })
      .finally(() => setLoading(false));
  };

  useEffect(() => {
    fetchBlocks();
  }, []);

  const handleCreateBlock = () => {
    const memberArray = members
      .split(",")
      .map((name) => ({ name: name.trim() }))
      .filter((m) => m.name);
    axios
      .post(
        "http://localhost:3000/blocks",
        { month, members: memberArray },
        { headers: { Authorization: `Bearer ${token}` } }
      )
      .then(() => {
        setOpen(false);
        setMonth("");
        setMembers("");
        fetchBlocks();
      })
      .catch((err) => console.error("Failed to create block", err));
  };

  const toggleLock = (month: string, locked: boolean) => {
    const url = `http://localhost:3000/blocks/${month}/${locked ? "unlock" : "lock"}`;
    axios
      .post(url, {}, { headers: { Authorization: `Bearer ${token}` } })
      .then(fetchBlocks)
      .catch((err) => console.error("Failed to toggle lock", err));
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
                    <Tooltip title={block.locked ? "Mở khóa tháng này" : "Khóa tháng này"}>
                      <IconButton
                        edge="end"
                        onClick={(e) => {
                          e.stopPropagation();
                          toggleLock(block.month, block.locked);
                        }}
                      >
                        {block.locked ? (
                          <LockIcon color="error" />
                        ) : (
                          <LockOpenIcon color="success" />
                        )}
                      </IconButton>
                    </Tooltip>
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
