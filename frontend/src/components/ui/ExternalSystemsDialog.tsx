import {
  CloseButton,
  Dialog,
  HStack,
  Portal,
  VStack,
  Icon,
  Text,
  Spinner,
  Box,
  Stack,
  Button,
  Input,
} from "@chakra-ui/react";
import { FiPlusSquare } from "react-icons/fi";
import {
  fetchExternalSystems,
  deleteExternalSystem,
  addExternalSystem,
} from "@/api/externalsystems";
import { useQuery, useQueryClient, useMutation } from "@tanstack/react-query";
import { toaster, Toaster } from "@/components/ui/toaster";
import { useState, useRef } from "react";

export default function ExternalSystemsDialog() {
  const [loadingIndex, setLoadingIndex] = useState<number | null>(null);
  const [externalName, setExternalName] = useState("");
  const { data, error, isLoading } = useQuery({
    queryKey: ["externalSystems"],
    queryFn: fetchExternalSystems,
  });

  const queryClient = useQueryClient();
  const ref = useRef<HTMLInputElement>(null);
  const deleteMutation = useMutation({
    mutationFn: (systemName: string) => deleteExternalSystem(systemName),
    onSuccess: () => {
      toaster.create({
        title: "Successful deletion!",
        description: "The system you have selected has been deleted",
        type: "success",
        duration: 2000,
      });

      queryClient.invalidateQueries({ queryKey: ["externalSystems"] });
    },
    onError: () => {
      toaster.create({
        title: "Failed to delete!",
        description: "The system you have selected has not been deleted",
        type: "error",
        duration: 2000,
      });
    },
    onSettled: () => {
      setLoadingIndex(null);
    },
  });

  const addMutation = useMutation({
    mutationFn: (systemName: string) => addExternalSystem(systemName),
    onSuccess: () => {
      toaster.create({
        title: "System added!",
        description: "The external system was successfully added.",
        type: "success",
        duration: 2000,
      });
      setExternalName("");
      queryClient.invalidateQueries({ queryKey: ["externalSystems"] });
    },
    onError: () => {
      toaster.create({
        title: "Failed to add system!",
        description: "An error occurred while adding the external system.",
        type: "error",
        duration: 2000,
      });
    },
  });

  const handleDelete = (name: string, index: number) => {
    setLoadingIndex(index);
    deleteMutation.mutate(name);
  };

  console.log(addMutation.data, externalName);

  return (
    <VStack>
      <Dialog.Root
        placement="center"
        motionPreset="slide-in-bottom"
        size="xl"
        initialFocusEl={() => ref.current}
      >
        <Dialog.Trigger asChild>
          <HStack
            gap={1}
            color="gray.600"
            _hover={{ color: "blue.500" }}
            cursor="pointer"
          >
            <Icon as={FiPlusSquare} />
            <Text fontSize="sm">External Systems</Text>
          </HStack>
        </Dialog.Trigger>
        <Portal>
          <Dialog.Backdrop />
          <Dialog.Positioner>
            <Dialog.Content>
              <Dialog.Header
                display="flex"
                alignItems="center"
                justifyContent="space-between"
              >
                <Dialog.Title>External systems</Dialog.Title>
                <Dialog.CloseTrigger asChild>
                  <CloseButton size="sm" />
                </Dialog.CloseTrigger>
              </Dialog.Header>

              <Dialog.Body pb="8">
              <Box pb="8">

              
              <HStack>

                
                  <Input
                    ref={ref}
                    value={externalName}
                    onChange={(e) => setExternalName(e.target.value)}
                    placeholder="Enter the name of an external system"
                    onKeyDown={(e) => {
                      if (e.key === "Enter" && externalName.trim() !== "") {
                        addMutation.mutate(externalName.trim());
                      }
                    }}
                  />
                  <Button
                    onClick={() => {
                      if (externalName.trim() !== "") {
                        addMutation.mutate(externalName.trim());
                      }
                    }}
                    disabled={externalName.trim() === ""}
                  >
                    Add external system
                  </Button>
                </HStack>
                {addMutation.data && (
                  <Box mt={4} p={4} borderWidth="1px" borderRadius="md">
                    <Text fontWeight="bold">Access Token for new system:</Text>
                    <Text color="green">{addMutation.data?.access_key}</Text>
                  </Box>
                )}
                </Box>
                {isLoading && (
                  <Box textAlign="center" py={4}>
                    <Spinner />
                    <Text mt={2}>Loading systems...</Text>
                  </Box>
                )}

                {error && (
                  <Box color="red.500" py={4}>
                    Failed to load systems.
                  </Box>
                )}
                <Stack gap={3} mb={6}>
                  {data?.slice().map((name: string, index: number) => (
                    <Box
                      key={index}
                      p={3}
                      borderWidth="1px"
                      borderRadius="lg"
                      display="flex"
                      alignItems="center"
                      justifyContent="space-between"
                    >
                      <Text fontWeight="medium">{name}</Text>
                      <HStack gap={3}>
                        {/* <Button
                          size="sm"
                          onClick={() => handleDelete(name, index)}
                          loading={loadingIndex === index}
                        >
                          Authenticate
                        </Button> */}

                        <Button
                          size="sm"
                          colorScheme="red"
                          onClick={() => handleDelete(name, index)}
                          loading={loadingIndex === index}
                        >
                          Delete
                        </Button>
                      </HStack>
                    </Box>
                  ))}
                </Stack>
                
              </Dialog.Body>
            </Dialog.Content>
          </Dialog.Positioner>
        </Portal>
      </Dialog.Root>
      <Toaster />
    </VStack>
  );
}
