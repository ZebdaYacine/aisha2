import type { Locale } from "./i18n";

export type CommonCopy = {
  loading: string;
  close: string;
  details: string;
  status: string;
  actions: string;
  previous: string;
  next: string;
  page: string;
  of: string;
  noResults: string;
  noMedia: string;
  save: string;
  update: string;
  delete: string;
  cancel: string;
  active: string;
  inactive: string;
  default: string;
  product: string;
  artisan: string;
  workshop: string;
  shop: string;
  reason: string;
  type: string;
  name: string;
  description: string;
  location: string;
  phone: string;
  search: string;
  openMedia: string;
  upload: string;
  refresh: string;
  yes: string;
  no: string;
  artisanSpecific: string;
  standardTraditional: string;
  chooseCategory: string;
};

const copies: Record<Locale, CommonCopy> = {
  en: { loading: "Loading…", close: "Close", details: "Details", status: "Status", actions: "Actions", previous: "Previous", next: "Next", page: "Page", of: "of", noResults: "No results found.", noMedia: "No media found.", save: "Save", update: "Update", delete: "Delete", cancel: "Cancel", active: "Active", inactive: "Inactive", default: "Default", product: "Product", artisan: "Artisan", workshop: "Workshop", shop: "Shop", reason: "Reason", type: "Type", name: "Name", description: "Description", location: "Location", phone: "Phone", search: "Search", openMedia: "Open media", upload: "Upload", refresh: "Refreshing…", yes: "Yes", no: "No", artisanSpecific: "Artisan specific", standardTraditional: "Standard traditional", chooseCategory: "Choose a category" },
  fr: { loading: "Chargement…", close: "Fermer", details: "Détails", status: "Statut", actions: "Actions", previous: "Précédent", next: "Suivant", page: "Page", of: "sur", noResults: "Aucun résultat trouvé.", noMedia: "Aucun média trouvé.", save: "Enregistrer", update: "Modifier", delete: "Supprimer", cancel: "Annuler", active: "Actif", inactive: "Inactif", default: "Par défaut", product: "Produit", artisan: "Artisan", workshop: "Atelier", shop: "Boutique", reason: "Motif", type: "Type", name: "Nom", description: "Description", location: "Adresse", phone: "Téléphone", search: "Rechercher", openMedia: "Ouvrir le média", upload: "Importer", refresh: "Actualisation…", yes: "Oui", no: "Non", artisanSpecific: "Spécifique à l’artisan", standardTraditional: "Traditionnel standard", chooseCategory: "Choisir une catégorie" },
  ar: { loading: "جارٍ التحميل…", close: "إغلاق", details: "التفاصيل", status: "الحالة", actions: "الإجراءات", previous: "السابق", next: "التالي", page: "الصفحة", of: "من", noResults: "لم يتم العثور على نتائج.", noMedia: "لا توجد وسائط.", save: "حفظ", update: "تحديث", delete: "حذف", cancel: "إلغاء", active: "نشط", inactive: "غير نشط", default: "افتراضي", product: "المنتج", artisan: "الحرفي", workshop: "الورشة", shop: "المتجر", reason: "السبب", type: "النوع", name: "الاسم", description: "الوصف", location: "الموقع", phone: "الهاتف", search: "بحث", openMedia: "فتح الوسائط", upload: "رفع", refresh: "جارٍ التحديث…", yes: "نعم", no: "لا", artisanSpecific: "خاص بالحرفي", standardTraditional: "تقليدي قياسي", chooseCategory: "اختر فئة" },
  es: { loading: "Cargando…", close: "Cerrar", details: "Detalles", status: "Estado", actions: "Acciones", previous: "Anterior", next: "Siguiente", page: "Página", of: "de", noResults: "No se encontraron resultados.", noMedia: "No se encontraron medios.", save: "Guardar", update: "Actualizar", delete: "Eliminar", cancel: "Cancelar", active: "Activo", inactive: "Inactivo", default: "Predeterminado", product: "Producto", artisan: "Artesano", workshop: "Taller", shop: "Tienda", reason: "Motivo", type: "Tipo", name: "Nombre", description: "Descripción", location: "Ubicación", phone: "Teléfono", search: "Buscar", openMedia: "Abrir medio", upload: "Subir", refresh: "Actualizando…", yes: "Sí", no: "No", artisanSpecific: "Específico del artesano", standardTraditional: "Tradicional estándar", chooseCategory: "Elegir una categoría" },
};

export function commonCopy(locale: Locale) {
  return copies[locale];
}

const statuses: Record<Locale, Record<string, string>> = {
  en: {},
  fr: { ACTIVE: "Actif", APPROVED: "Approuvé", VERIFIED: "Vérifié", INSPECTED: "Inspecté", DELIVERED: "Livré", PAID: "Payé", PROCESSING: "En traitement", PENDING: "En attente", PENDING_REVIEW: "En attente d’examen", RECEIVED_PENDING_INSPECTION: "Réception à inspecter", CHANGES_REQUESTED: "Modifications demandées", SUSPENDED: "Suspendu", REJECTED: "Rejeté", DISABLED: "Désactivé", CANCELLED: "Annulé", RETURNED: "Retourné", DAMAGED: "Endommagé", QUARANTINED: "En quarantaine", ARCHIVED: "Archivé", INACTIVE: "Inactif", AVAILABLE: "Disponible", OUT_OF_STOCK: "Rupture de stock", NOT_STARTED: "Non commencé", DRAFT: "Brouillon" },
  ar: { ACTIVE: "نشط", APPROVED: "معتمد", VERIFIED: "موثّق", INSPECTED: "تم التفتيش", DELIVERED: "تم التوصيل", PAID: "مدفوع", PROCESSING: "قيد المعالجة", PENDING: "قيد الانتظار", PENDING_REVIEW: "قيد المراجعة", RECEIVED_PENDING_INSPECTION: "مستلم بانتظار التفتيش", CHANGES_REQUESTED: "مطلوب تعديل", SUSPENDED: "معلّق", REJECTED: "مرفوض", DISABLED: "معطّل", CANCELLED: "ملغى", RETURNED: "مرتجع", DAMAGED: "تالف", QUARANTINED: "في الحجر", ARCHIVED: "مؤرشف", INACTIVE: "غير نشط", AVAILABLE: "متاح", OUT_OF_STOCK: "نفد المخزون", NOT_STARTED: "لم يبدأ", DRAFT: "مسودة" },
  es: { ACTIVE: "Activo", APPROVED: "Aprobado", VERIFIED: "Verificado", INSPECTED: "Inspeccionado", DELIVERED: "Entregado", PAID: "Pagado", PROCESSING: "En proceso", PENDING: "Pendiente", PENDING_REVIEW: "Pendiente de revisión", RECEIVED_PENDING_INSPECTION: "Recibido, pendiente de inspección", CHANGES_REQUESTED: "Cambios solicitados", SUSPENDED: "Suspendido", REJECTED: "Rechazado", DISABLED: "Desactivado", CANCELLED: "Cancelado", RETURNED: "Devuelto", DAMAGED: "Dañado", QUARANTINED: "En cuarentena", ARCHIVED: "Archivado", INACTIVE: "Inactivo", AVAILABLE: "Disponible", OUT_OF_STOCK: "Agotado", NOT_STARTED: "No iniciado", DRAFT: "Borrador" },
};

export function statusLabel(status: string, locale: Locale) {
  return statuses[locale][status] ?? status.replaceAll("_", " ");
}
