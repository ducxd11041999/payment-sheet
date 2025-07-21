import { useCallback, useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import {
  Typography,
  Table,
  TableHead,
  TableRow,
  TableCell,
  TableBody,
  CircularProgress,
  Card,
  CardContent,
  Box,
  Button,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  IconButton,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  TextField,
  Collapse,
} from "@mui/material";
import Grid from "@mui/material/Grid";
import DeleteIcon from "@mui/icons-material/Delete";
import EditIcon from "@mui/icons-material/Edit";
import ArrowBackIcon from "@mui/icons-material/ArrowBack";
import {
  getTransactions,
  getMembers,
  addTransaction,
  deleteTransaction,
  updateTransaction,
} from "../api/api";

interface Transaction {
  id: string;
  description: string;
  amount: number;
  payer: string;
  created_at: string;
  ratios: Record<string, number>;
  details: Record<string, number>;
}

interface Member {
  id: string;
  name: string;
  ratio: number;
  debt: number;
}

export default function BlockDetail() {
  const { month } = useParams();
  const navigate = useNavigate();
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [members, setMembers] = useState<Member[]>([]);
  const [loading, setLoading] = useState(true);
  const [openAdd, setOpenAdd] = useState(false);
  const [newDesc, setNewDesc] = useState("");
  const [newAmount, setNewAmount] = useState("");
  const [newPayer, setNewPayer] = useState("");
  const [newRatios, setNewRatios] = useState<Record<string, number>>({});
  const [expandedRow, setExpandedRow] = useState<string | null>(null);
  const [editTransaction, setEditTransaction] = useState<Transaction | null>(null);
  const token = localStorage.getItem("token");

  const fetchData = useCallback(async () => {
    try {
      const [transRes, membersRes] = await Promise.all([
        getTransactions(month!),
        getMembers(month!),
      ]);
      setTransactions(transRes.data || []);
      setMembers(membersRes.data || []);
    } catch (err) {
      console.error("Failed to load block details", err);
      setTransactions([]);
      setMembers([]);
    } finally {
      setLoading(false);
    }
  }, [month, token]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  const memberMap = members.reduce<Record<string, string>>((acc, m) => {
    acc[m.id] = m.name;
    return acc;
  }, {});

  const handleDeleteTransaction = (id: string) => {
    deleteTransaction(id)
      .then(fetchData)
      .catch((err) => console.error("Failed to delete transaction", err));
  };

  const openAddDialog = () => {
    const initialRatios = members.reduce<Record<string, number>>((acc, m) => {
      acc[m.id] = 1;
      return acc;
    }, {});
    setNewDesc("");
    setNewAmount("");
    setNewPayer("");
    setNewRatios(initialRatios);
    setEditTransaction(null);
    setOpenAdd(true);
  };

  const openEditDialog = (transaction: Transaction) => {
    setNewDesc(transaction.description);
    setNewAmount(transaction.amount.toString());
    setNewPayer(transaction.payer);
    setNewRatios(transaction.ratios);
    setEditTransaction(transaction);
    setOpenAdd(true);
  };

  const handleAddTransaction = () => {
    if (editTransaction) {
      updateTransaction(editTransaction.id, {
        description: newDesc,
        amount: parseFloat(newAmount),
        payer: newPayer,
        ratios: newRatios,
      })
        .then(() => {
          setOpenAdd(false);
          setEditTransaction(null);
          fetchData();
        })
        .catch((err) => console.error("Failed to update transaction", err));
    } else {
      addTransaction(month!, {
        description: newDesc,
        amount: parseFloat(newAmount),
        payer: newPayer,
        ratios: newRatios,
      })
        .then(() => {
          setOpenAdd(false);
          fetchData();
        })
        .catch((err) => console.error("Failed to add transaction", err));
    }
  };

  const toggleRow = (id: string) => {
    setExpandedRow(expandedRow === id ? null : id);
  };

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" height="100vh">
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box sx={{ padding: 3, backgroundColor: "#f5f6fa", minHeight: "100vh" }}>
      {/* Header */}
      <Box display="flex" alignItems="center" justifyContent="space-between" mb={2}>
        <IconButton onClick={() => navigate("/home")}>
          <ArrowBackIcon />
        </IconButton>
        <Typography variant="h5" align="center" sx={{ flexGrow: 1, fontWeight: "bold" }}>
          Chi tiết tháng {month}
        </Typography>
        <Box width="48px" />
      </Box>

      <Grid container spacing={3}>
        {/* Bảng chi tiêu */}
        <Grid size={{ xs: 12, md: 8 }}>
          <Card sx={{ boxShadow: 3 }}>
            <CardContent>
              <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
                <Typography variant="h6" sx={{ fontWeight: "bold" }}>
                  Danh sách chi tiêu
                </Typography>
                <Button variant="contained" color="primary" onClick={openAddDialog}>
                  Thêm chi tiêu
                </Button>
              </Box>
              <Table sx={{ width: "100%" }}>
                <TableHead sx={{ backgroundColor: "#e0e0e0" }}>
                  <TableRow>
                    <TableCell sx={{ fontWeight: "bold" }}>#</TableCell>
                    <TableCell sx={{ fontWeight: "bold" }}>Mô tả</TableCell>
                    <TableCell sx={{ fontWeight: "bold" }} align="right">
                      Số tiền
                    </TableCell>
                    <TableCell sx={{ fontWeight: "bold" }} align="right">
                      Ngày
                    </TableCell>
                    <TableCell sx={{ fontWeight: "bold" }} align="center">
                      Hành động
                    </TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {transactions.length > 0 ? (
                    transactions.map((t, index) => (
                      <>
                        <TableRow
                          key={t.id}
                          hover
                          onClick={() => toggleRow(t.id)}
                          style={{ cursor: "pointer" }}
                        >
                          <TableCell>{index + 1}</TableCell>
                          <TableCell>{t.description}</TableCell>
                          <TableCell align="right" sx={{ color: "green", fontWeight: "bold" }}>
                            {t.amount.toLocaleString()} ₫
                          </TableCell>
                          <TableCell align="right">
                            {new Date(t.created_at).toLocaleDateString()}
                          </TableCell>
                          <TableCell align="center">
                            <IconButton
                              onClick={(e) => {
                                e.stopPropagation();
                                openEditDialog(t);
                              }}
                              color="primary"
                            >
                              <EditIcon />
                            </IconButton>
                            <IconButton
                              onClick={(e) => {
                                e.stopPropagation();
                                handleDeleteTransaction(t.id);
                              }}
                              color="error"
                            >
                              <DeleteIcon />
                            </IconButton>
                          </TableCell>
                        </TableRow>
                        <TableRow>
                          <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={5}>
                            <Collapse in={expandedRow === t.id} timeout="auto" unmountOnExit>
                              <Box margin={2}>
                                <Typography variant="subtitle1" sx={{ fontWeight: "bold" }}>
                                  Người trả:{" "}
                                  <span style={{ color: "#1976d2" }}>
                                    {memberMap[t.payer] || t.payer}
                                  </span>
                                </Typography>
                                <Typography variant="body2" sx={{ mt: 1 }}>
                                  Tỉ lệ chia:{" "}
                                  {Object.entries(t.ratios).map(([memberId, value], idx) => (
                                    <span key={memberId}>
                                      <span style={{ fontWeight: "bold", color: "#1976d2" }}>
                                        {memberMap[memberId]}
                                      </span>
                                      : <span style={{ color: "#555" }}>{value}</span>
                                      {idx < Object.entries(t.ratios).length - 1 && ", "}
                                    </span>
                                  ))}
                                </Typography>
                                {t.details && (
                                  <Typography variant="body2" sx={{ mt: 1 }}>
                                    Số tiền mỗi người:{" "}
                                    {Object.entries(t.details).map(([memberId, amount], idx) => (
                                      <span key={memberId}>
                                        <span style={{ fontWeight: "bold", color: "#1976d2" }}>
                                          {memberMap[memberId]}
                                        </span>
                                        :{" "}
                                        <span style={{ fontWeight: "bold", color: "green" }}>
                                          {amount.toLocaleString()} ₫
                                        </span>
                                        {idx < Object.entries(t.details).length - 1 && ", "}
                                      </span>
                                    ))}
                                  </Typography>
                                )}
                              </Box>
                            </Collapse>
                          </TableCell>
                        </TableRow>
                      </>
                    ))
                  ) : (
                    <TableRow>
                      <TableCell colSpan={5} align="center">
                        Không có chi tiêu nào
                      </TableCell>
                    </TableRow>
                  )}
                </TableBody>
              </Table>
            </CardContent>
          </Card>
        </Grid>

        {/* Bảng tổng kết thành viên */}
        <Grid size={{ xs: 12, md: 4 }}>
          <Card sx={{ boxShadow: 3 }}>
            <CardContent>
              <Typography variant="h6" gutterBottom align="center" sx={{ fontWeight: "bold" }}>
                Tổng kết thành viên
              </Typography>

              <Typography variant="subtitle1" align="center" sx={{ mb: 2, fontWeight: "bold", color: "green" }}>
                Tổng chi tiêu phòng: {transactions.reduce((sum, t) => sum + t.amount, 0).toLocaleString()} ₫
              </Typography>

              <Table>
                <TableHead sx={{ backgroundColor: "#e0e0e0" }}>
                  <TableRow>
                    <TableCell sx={{ fontWeight: "bold" }}>Tên</TableCell>
                    <TableCell sx={{ fontWeight: "bold" }} align="right">
                      Đã trả
                    </TableCell>
                    <TableCell sx={{ fontWeight: "bold" }} align="right">
                      Nợ
                    </TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {members.map((m) => {
                    const paid = transactions
                      .filter((t) => t.payer === m.id)
                      .reduce((sum, t) => sum + t.amount, 0);
                    return (
                      <TableRow key={m.id}>
                        <TableCell sx={{ fontWeight: "bold", color: "#1976d2" }}>{m.name}</TableCell>
                        <TableCell align="right" sx={{ color: "blue", fontWeight: "bold" }}>
                          {paid.toLocaleString()} ₫
                        </TableCell>
                        <TableCell align="right" sx={{ color: m.debt < 0 ? "red" : "green", fontWeight: "bold" }}>
                          {m.debt.toLocaleString()} ₫
                        </TableCell>
                      </TableRow>
                    );
                  })}
                </TableBody>
              </Table>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      {/* Dialog thêm/sửa chi tiêu */}
      <Dialog open={openAdd} onClose={() => setOpenAdd(false)} maxWidth="sm" fullWidth>
        <DialogTitle>{editTransaction ? "Chỉnh sửa chi tiêu" : "Thêm chi tiêu mới"}</DialogTitle>
        <DialogContent>
          <TextField
            label="Mô tả"
            fullWidth
            margin="normal"
            value={newDesc}
            onChange={(e) => setNewDesc(e.target.value)}
          />
          <TextField
            label="Số tiền"
            type="number"
            fullWidth
            margin="normal"
            value={newAmount}
            onChange={(e) => setNewAmount(e.target.value)}
          />
          <FormControl fullWidth margin="normal">
            <InputLabel id="payer-label">Người trả</InputLabel>
            <Select
                labelId="payer-label"
                value={newPayer}
                onChange={(e) => setNewPayer(e.target.value)}
                label="Người trả"
            >
              {members.map((m) => (
                  <MenuItem key={m.id} value={m.id}>
                    {m.name}
                  </MenuItem>
              ))}
            </Select>
          </FormControl>
          <Typography variant="subtitle1" sx={{ mt: 2, mb: 1 }}>
            Tỉ lệ cho từng thành viên
          </Typography>
          {members.map((m) => (
            <TextField
              key={m.id}
              label={`${m.name}`}
              type="number"
              fullWidth
              margin="dense"
              value={newRatios[m.id] ?? 0}
              onChange={(e) =>
                setNewRatios((prev) => ({
                  ...prev,
                  [m.id]: parseFloat(e.target.value) < 0? 0: parseFloat(e.target.value),
                }))
              }
            />
          ))}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenAdd(false)}>Hủy</Button>
          <Button onClick={handleAddTransaction} variant="contained">
            {editTransaction ? "Cập nhật" : "Thêm"}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}
