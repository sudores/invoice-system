<template>
  <div class="layout">
    <!-- Content -->
    <main class="content">
      <!-- Summary cards -->
      <div class="cards">
        <div class="card">
          <div class="card-title">Total Unpaid</div>
          <div class="card-value">{{ summary.unpaid }}</div>
        </div>

        <div class="card">
          <div class="card-title">Total Paid</div>
          <div class="card-value">{{ summary.paid }}</div>
        </div>

        <div class="card">
          <div class="card-title">Due Soon</div>
          <div class="card-value">{{ summary.dueSoon }}</div>
        </div>
      </div>

      <!-- Invoices List -->
      <div class="list-card">
        <h3>Your Invoices</h3>

        <div v-if="loading" class="loading">Loading…</div>
        <div v-if="error" class="error">{{ error }}</div>

        <div v-for="inv in invoices" :key="inv.id" class="invoice-row">
          <div class="left">
            <div class="inv-title">{{ inv.title }}</div>
            <div class="inv-meta">
              {{ inv.amount }} — {{ inv.status.toUpperCase() }}
            </div>
          </div>

          <div class="right">
            <button
              v-if="inv.status !== 'PAID'"
              class="btn small"
              @click="markPaid(inv.id)"
            >Mark Paid</button>

            <button
              class="btn small gray"
              @click="edit(inv.id)"
            >Edit</button>
          </div>
        </div>

        <div v-if="!loading && invoices.length === 0" class="empty">
          No invoices found.
        </div>
      </div>

      <button class="btn primary full" @click="create">Create Invoice</button>
    </main>
  </div>
</template>

<script>
export default {
  data() {
    return {
      loading: true,
      error: null,
      invoices: [],
      summary: {
        unpaid: 0,
        paid: 0,
        dueSoon: 0,
      },
    };
  },

  async mounted() {
    const token = localStorage.getItem("token");
    if (!token) {
      this.$router.push("/login");
      return;
    }

    try {
      const resp = await fetch("/api/invoices", {
        headers: { Authorization: `Bearer ${token}` },
      });

      if (!resp.ok) throw new Error(await resp.text());
      this.invoices = await resp.json();

      this.computeSummary();
    } catch (err) {
      this.error = err.message;
    } finally {
      this.loading = false;
    }
  },

  methods: {
    computeSummary() {
      let unpaid = 0;
      let paid = 0;
      let dueSoon = 0;

      const now = Date.now();
      const soon = 1000 * 60 * 60 * 24 * 7; // 7 days

      for (const inv of this.invoices) {
        if (inv.status === "PAID") paid += inv.amount;
        else unpaid += inv.amount;

        if (inv.due_date && new Date(inv.due_date).getTime() - now < soon)
          dueSoon++;
      }

      this.summary = {
        unpaid: unpaid.toFixed(2),
        paid: paid.toFixed(2),
        dueSoon,
      };
    },

    async markPaid(id) {
      const token = localStorage.getItem("token");
      await fetch(`/api/invoices/${id}/mark-paid`, {
        method: "POST",
        headers: { Authorization: `Bearer ${token}` },
      });

      this.invoices = this.invoices.map((i) =>
        i.id === id ? { ...i, status: "PAID" } : i
      );

      this.computeSummary();
    },

    create() {
      this.$router.push("/create");
    },

    edit(id) {
      this.$router.push(`/edit/${id}`);
    },

    logout() {
      localStorage.removeItem("token");
      this.$router.push("/login");
    },
  },
};
</script>

<style scoped>
/* Global Layout */
.layout {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: #f5f7fa;
  font-family: system-ui, sans-serif;
}

/* Navigation */
.nav {
  height: 56px;
  background: #2c3e50;
  color: #fff;
  padding: 0 16px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.nav-title {
  font-size: 1.2rem;
  font-weight: 600;
}

.logout {
  background: #e74c3c;
  border: none;
  padding: 6px 12px;
  border-radius: 4px;
  color: white;
  cursor: pointer;
}
.logout:hover {
  background: #c0392b;
}

/* Main content */
.content {
  padding: 0;
  overflow-y: auto;
}

/* Summary Cards */
.cards {
  display: flex;
  gap: 12px;
  margin-bottom: 18px;
}

.card {
  flex: 1;
  background: white;
  padding: 16px;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.08);
}

.card-title {
  font-size: 0.85rem;
  color: #8b8b8b;
}

.card-value {
  margin-top: 8px;
  font-size: 1.4rem;
  font-weight: 600;
}

/* Invoices List */
.list-card {
  background: white;
  padding: 18px;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.08);
  margin-bottom: 18px;
}

.invoice-row {
  display: flex;
  justify-content: space-between;
  padding: 12px 0;
  border-bottom: 1px solid #e5e5e5;
}

.invoice-row:last-child {
  border-bottom: none;
}

.inv-title {
  font-weight: 600;
}

.inv-meta {
  font-size: 0.85rem;
  color: #777;
}

.empty {
  text-align: center;
  padding: 20px;
  color: #777;
}

/* Buttons */
.btn {
  border: none;
  padding: 8px 12px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.9rem;
}

.btn.primary {
  background: #3498db;
  color: white;
}

.btn.primary:hover {
  background: #2980b9;
}

.btn.gray {
  background: #bdc3c7;
  color: #2c3e50;
}

.btn.gray:hover {
  background: #a6acaf;
}

.btn.small {
  padding: 6px 10px;
  font-size: 0.8rem;
}

.btn.full {
  width: 100%;
  margin-top: 10px;
}

/* Loading & error */
.loading {
  color: #777;
}

.error {
  color: #c0392b;
  margin-bottom: 12px;
}

/* Mobile stacking */
@media (max-width: 600px) {
  .cards {
    flex-direction: column;
  }
}
</style>

