<template>
  <div class="flex flex-col h-screen bg-gray-100 font-sans">
    <!-- Main content -->
    <main class="flex-1 p-4 overflow-y-auto">
      <!-- Summary Cards -->
      <div class="flex flex-col sm:flex-row gap-3 mb-4">
        <div class="flex-1 bg-white p-4 rounded-lg shadow">
          <div class="text-gray-500 text-sm">Total Unpaid</div>
          <div class="mt-2 text-xl font-semibold">{{ summary.unpaid }}</div>
        </div>

        <div class="flex-1 bg-white p-4 rounded-lg shadow">
          <div class="text-gray-500 text-sm">Total Paid</div>
          <div class="mt-2 text-xl font-semibold">{{ summary.paid }}</div>
        </div>

        <div class="flex-1 bg-white p-4 rounded-lg shadow">
          <div class="text-gray-500 text-sm">Due Soon</div>
          <div class="mt-2 text-xl font-semibold">{{ summary.dueSoon }}</div>
        </div>
      </div>

      <!-- Invoices List -->
      <div class="bg-white p-4 rounded-lg shadow mb-4">
        <h3 class="text-lg font-semibold mb-2">Your Invoices</h3>

        <div v-if="loading" class="text-gray-500">Loading…</div>
        <div v-if="error" class="text-red-600 mb-3">{{ error }}</div>

        <div
          v-for="inv in invoices"
          :key="inv.id"
          class="flex justify-between py-3 border-b last:border-b-0"
        >
          <div>
            <div class="font-semibold">{{ inv.title }}</div>
            <div class="text-sm text-gray-600">
              {{ inv.amount }} — {{ inv.status.toUpperCase() }}
            </div>
          </div>

          <div class="flex gap-2">
            <button
              v-if="inv.status !== 'PAID'"
              @click="markPaid(inv.id)"
              class="px-2 py-1 text-sm rounded bg-blue-500 text-white hover:bg-blue-600"
            >
              Mark Paid
            </button>

            <button
              @click="edit(inv.id)"
              class="px-2 py-1 text-sm rounded bg-gray-300 text-gray-800 hover:bg-gray-400"
            >
              Edit
            </button>
          </div>
        </div>

        <div v-if="!loading && invoices.length === 0" class="text-center py-5 text-gray-600">
          No invoices found.
        </div>
      </div>

      <button
        @click="create"
        class="w-full px-4 py-2 text-white bg-blue-500 rounded hover:bg-blue-600"
      >
        Create Invoice
      </button>
    </main>
  </div>
</template>

<script>
import api from '../../api/client'
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
      const resp = await api.get("/api/v1/invoice/list/sent", {
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
